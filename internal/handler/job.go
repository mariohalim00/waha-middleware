package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/service"
	"waha-job-processing/internal/util"
)

var TEXT_TEMPLATE = os.Getenv("TEXT_BLAST_TEMPLATE")

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

	go procesJobBackground(jobList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{
	"status": "accepted",
	"message": "Job is being processed in the background. Results will be sent to your webhook."
}`))
}

func procesJobBackground(jobList []models.Job) {
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
		if idx < length {
			time.Sleep(util.GenerateRandomDuration(30))
		}
		//	return list of failed jobs
	}

	body, err := json.Marshal(failedJobs)

	if err != nil {
		fmt.Println("Error converting to JSON: ", err)
	}

	util.Post(body, os.Getenv("BLASTER_WEBHOOK_URL"))

}
