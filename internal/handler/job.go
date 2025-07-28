package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/service"
	"waha-job-processing/internal/util"
	"waha-job-processing/internal/util/httpHelper"

	"github.com/joho/godotenv"
)

var TEXT_TEMPLATE string
var PROMO_CODE string
var WEBFORM_URL string

func init() {
	godotenv.Load()
	TEXT_TEMPLATE = os.Getenv("TEXT_BLAST_TEMPLATE")
	PROMO_CODE = os.Getenv("PROMO_CODE")
	WEBFORM_URL = os.Getenv("BASE_WEBFORM_URL")
}

var jobProcessing = make(map[string]bool)

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

	// Create a unique key for this job batch (could be improved with a real idempotency key)
	jobKey := fmt.Sprintf("%v", jobList)
	if jobProcessing[jobKey] {
		httpHelper.ReturnHttpError(w, "Job is already being processed", http.StatusConflict)
		return
	}
	jobProcessing[jobKey] = true

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in job processor: %v", r)
			}
			delete(jobProcessing, jobKey)
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
		log.Println("Processing job for:", job.Customer.FormattedPhoneNumber, "Session:", job.Pic.Session)
		err := service.StartTyping(job.Pic.Session, job.Customer.FormattedPhoneNumber)
		if err != nil {
			// httpHelper.ReturnHttpError(w, "Error sending typing event", http.StatusConflict)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
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
				Voucher:        job.Customer.Voucher,
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
				Voucher:        job.Customer.Voucher,
			})
			continue
		}
		msg := strings.ReplaceAll(TEXT_TEMPLATE, "{{name}}", job.Customer.Name)
		msg = strings.ReplaceAll(msg, `\n`, "\n")
		msg = msg + "\n" + url
		err = service.SendMessage(job.Pic.Session, job.Customer.FormattedPhoneNumber, msg)
		if err != nil {
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
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

	err = callWebhookPostJob(body)
	if err != nil {
		log.Printf("Error when calling webhook: %+v\n", err)
	}

}

func callWebhookPostJob(body []byte) error {
	var err error
	MAX_RETRY := 3
	for i := range MAX_RETRY {
		err = httpHelper.Post(body, os.Getenv("BLASTER_WEBHOOK_URL"), httpHelper.HttpHeader{
			"Content-Type": "application/json",
		})
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

func generateWebFormUrl(jobData models.Job) (string, error) {
	currDate := time.Now()
	log.Printf("generating job for current voucher: %+v, job: %+v", jobData.Customer.Voucher, jobData)

	assignedVoucher, err := service.GetVoucherByName(jobData.Customer.Voucher)
	if err != nil {
		log.Printf("Error retrieving voucher: %v", err)
		return "", fmt.Errorf("failed to retrieve voucher: %w", err)
	}

	duration := time.Duration(assignedVoucher.PromoDurationHours) * time.Hour
	expiryDate := currDate.Add(duration)
	var promoToken = models.PromoToken{
		UserName:  jobData.Customer.Name,
		ExpiresAt: expiryDate,
		PromoCode: jobData.Customer.Voucher,
	}

	signature, err := service.BuildSignature(promoToken, string(rune(currDate.Unix())))

	if err != nil {
		return "", err
	}

	signedUrl := WEBFORM_URL + "/web-form?data=" + signature
	_, err = service.CreateTrackedPromo(signature, jobData.Customer.Name, jobData.Customer.Voucher, expiryDate)

	if err != nil {
		log.Printf("Error inserting %v", err)
	}

	return signedUrl, nil
}
