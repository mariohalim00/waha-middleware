package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/service"
	"waha-job-processing/internal/util"
)

var TEXT_TEMPLATE string

func init() {
	TEXT_TEMPLATE = os.Getenv("TEXT_BLAST_TEMPLATE")
	if TEXT_TEMPLATE == "" {
		panic("TEXT_BLAST_TEMPLATE environment variable is not set")
	}
}

func ProcessJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jobList []models.Job
	err := json.NewDecoder(r.Body).Decode(&jobList)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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
		err := service.StartTyping(job.Pic.Session, job.Customer.ChatId)
		if err != nil {
			// http.Error(w, "Error sending typing event", http.StatusConflict)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		time.Sleep(util.GenerateRandomDuration(30))

		//	stop typing
		err = service.StopTyping(job.Pic.Session, job.Customer.ChatId)
		if err != nil {
			// http.Error(w, "Error stopping typing event", http.StatusConflict)
			failedJobs = append(failedJobs, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		time.Sleep(util.GenerateRandomDuration(30))

		//	send message
		msg := strings.ReplaceAll(TEXT_TEMPLATE, "{{name}}", job.Customer.Name)
		msg = strings.ReplaceAll(msg, `\n`, "\n")
		err = service.SendMessage(job.Pic.Session, job.Customer.ChatId, msg)
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
		log.Printf("Error when callling webhook: %+v\n", err)

	}

}

func callWebhook(body []byte) error {
	var err error
	MAX_RETRY := 3
	for i := range MAX_RETRY {
		err = util.Post(body, os.Getenv("BLASTER_WEBHOOK_URL"))
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
