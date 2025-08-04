package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"waha-job-processing/internal/database/repository"
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

	// Create a unique key for this job batch based on phone numbers and session
	var phoneNumbers []string
	for _, job := range jobList {
		phoneNumbers = append(phoneNumbers, job.Customer.FormattedPhoneNumber)
	}
	jobKey := strings.Join(phoneNumbers, ",")

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
			// Always clean up the job key when done
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
		custdata, err := service.GetUserDataByUsername(job.Customer.Username)
		if err != nil {
			log.Println("Error fetching user data:", err)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
				Username:       job.Customer.Username,
			})
			continue
		}

		assignedVoucher, err := service.GetVoucherByName(job.Customer.Voucher)
		if err != nil {
			log.Printf("Error retrieving voucher: %v", err)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
				Username:       job.Customer.Username,
			})
			continue
		}

		//	start typing
		log.Println("Processing job for:", job.Customer.FormattedPhoneNumber, "Session:", job.Pic.Session)
		err = service.StartTyping(job.Pic.Session, job.Customer.FormattedPhoneNumber)
		if err != nil {
			// httpHelper.ReturnHttpError(w, "Error sending typing event", http.StatusConflict)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
				Username:       job.Customer.Username,
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
				Username:       job.Customer.Username,
			})
			continue
		}
		time.Sleep(util.GenerateRandomDuration(30))

		//	send message
		url, err := generateWebFormUrl(job, assignedVoucher, custdata.Userid)
		if err != nil {
			log.Printf("failed to generate url, skipping")
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
				Username:       job.Customer.Username,
			})
			continue
		}

		var msg_template string
		if assignedVoucher.PromoTextTemplate.Valid && assignedVoucher.PromoTextTemplate.String != "" {
			msg_template = assignedVoucher.PromoTextTemplate.String
		} else {
			msg_template = TEXT_TEMPLATE
		}

		var name string
		if custdata.Name != "" {
			name = custdata.Name
		} else {
			name = job.Customer.Name
		}

		msg := strings.ReplaceAll(msg_template, "{{name}}", name)
		msg = strings.ReplaceAll(msg, `\n`, "\n")
		msg = msg + "\n" + url

		err = service.SendMessage(job.Pic.Session, job.Customer.FormattedPhoneNumber, msg)
		if err != nil {
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
				Voucher:        job.Customer.Voucher,
				Username:       job.Customer.Username,
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

func generateWebFormUrl(jobData models.Job, assignedVoucher *repository.Voucher, userID string) (string, error) {
	currDate := time.Now()
	log.Printf("generating job for current voucher: %+v, job: %+v", jobData.Customer.Voucher, jobData)

	duration := time.Duration(assignedVoucher.PromoDurationHours) * time.Hour
	expiryDate := currDate.Add(duration)
	var promoToken = models.PromoToken{
		UserName:  jobData.Customer.Username,
		ExpiresAt: expiryDate,
		PromoCode: jobData.Customer.Voucher,
	}

	signature, err := service.BuildSignature(promoToken, string(rune(currDate.Unix())))

	if err != nil {
		return "", err
	}

	signedUrl := WEBFORM_URL + "/web-form?data=" + signature
	_, err = service.CreateTrackedPromo(signature, jobData.Customer.Username, jobData.Customer.Voucher, expiryDate, userID)

	if err != nil {
		log.Printf("Error inserting %v", err)
	}

	return signedUrl, nil
}
