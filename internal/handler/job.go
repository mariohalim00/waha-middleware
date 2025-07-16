package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"waha-job-processing/internal/database/db"
	"waha-job-processing/internal/database/repository"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/service"
	"waha-job-processing/internal/util"
	"waha-job-processing/internal/util/httpHelper"

	"github.com/jackc/pgx/v5/pgtype"
)

var TEXT_TEMPLATE string
var PROMO_CODE string
var WEBFORM_URL string

func init() {
	TEXT_TEMPLATE = os.Getenv("TEXT_BLAST_TEMPLATE")
	PROMO_CODE = os.Getenv("PROMO_CODE")
	WEBFORM_URL = os.Getenv("BASE_WEBFORM_URL")
}

func ProcessJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jobList []models.Job
	err := json.NewDecoder(r.Body).Decode(&jobList)

	if err != nil {
		httpHelper.ReturnHttpError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in job processor: %v", r)
			}
		}()
		processJobBackground(jobList)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{
	"status": "accepted",
	"message": "Job is being processed in the background. Results will be sent to your webhook."
}`))
}

func processJobBackground(jobList []models.Job) {
	var failedJobs []models.JobResponse

	length := len(jobList)
	for idx, job := range jobList {
		//	start typing
		fmt.Println("Processing job for:", job.Customer.FormattedPhoneNumber, "Session:", job.Pic.Session)
		err := service.StartTyping(job.Pic.Session, job.Customer.FormattedPhoneNumber)
		if err != nil {
			// httpHelper.ReturnHttpError(w, "Error sending typing event", http.StatusConflict)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		time.Sleep(util.GenerateRandomDuration(30))

		//	stop typing
		err = service.StopTyping(job.Pic.Session, job.Customer.FormattedPhoneNumber)
		if err != nil {
			// httpHelper.ReturnHttpError(w, "Error stopping typing event", http.StatusConflict)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		time.Sleep(util.GenerateRandomDuration(30))

		//	send message
		url, err := generateWebFormUrl(job)
		if err != nil {
			log.Printf("failed to generate url, skipping")
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		msg := strings.ReplaceAll(TEXT_TEMPLATE, "{{name}}", job.Customer.Name)
		msg = strings.ReplaceAll(msg, `\n`, "\n"+url)
		err = service.SendMessage(job.Pic.Session, job.Customer.FormattedPhoneNumber, msg)
		if err != nil {
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		// wait few secs
		if idx < length-1 {
			time.Sleep(util.GenerateRandomDuration(30))
		}
		//	return list of failed jobs
	}

	body, err := json.Marshal(failedJobs)

	if err != nil {
		log.Printf("Error converting failedJobs to JSON: %v. FailedJobs content: %+v\n", err, failedJobs)
	}

	err = callWebhook(body)
	if err != nil {
		log.Printf("Error when calling webhook: %+v\n", err)
	}

}

func callWebhook(body []byte) error {
	var err error
	MAX_RETRY := 3
	for i := range MAX_RETRY {
		err = httpHelper.Post(body, os.Getenv("BLASTER_WEBHOOK_URL"))
		if err == nil {
			break
		}
		if i < MAX_RETRY {
			log.Printf("Failed to call webhook, attempt %v/%v", i, MAX_RETRY)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func trackDistributedPromo(signature, userName string, expiryDate time.Time) (repository.PromoTracker, error) {
	ctx := context.Background()
	conn := db.New(ctx)
	defer conn.Close(ctx)

	param := repository.CreateTrackedPromoParams{
		HashedString: signature,
		ExpiredAt: pgtype.Timestamptz{
			Time:  expiryDate,
			Valid: true,
		},
		UserName: userName,
	}

	query := repository.New(conn)

	trackedPromo, err := query.CreateTrackedPromo(ctx, param)
	if err != nil {
		log.Printf("Failed to crate tracked promo record")
		return repository.PromoTracker{}, err
	}

	log.Printf("Created tracked promo %v\n", trackedPromo)
	return trackedPromo, nil
}

func generateWebFormUrl(jobData models.Job) (string, error) {
	currDate := time.Now()
	expiryDate := currDate.AddDate(0, 0, 7)
	var promoToken = models.PromoToken{
		UserName:  jobData.Customer.Name,
		ExpiresAt: expiryDate,
		PromoCode: PROMO_CODE,
	}

	signature, err := service.BuildSignature(promoToken, string(rune(currDate.Unix())))

	if err != nil {
		return "", err
	}

	signedUrl := WEBFORM_URL + "/web-form?data=" + signature
	_, err = trackDistributedPromo(signature, jobData.Customer.Name, expiryDate)

	if err != nil {
		log.Printf("Error inserting %v", err)
	}

	return signedUrl, nil
}
