package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/service"
	"waha-job-processing/internal/util"
)

var TEXT_TEMPLATE = os.Getenv("TEXT_BLAST_TEMPLATE")

func ProcessJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jobList []models.Job
	var failedJob []models.JobResponse
	err := json.NewDecoder(r.Body).Decode(&jobList)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	length := len(jobList)
	for idx, job := range jobList {
		//	start typing
		err = service.StartTyping(job.Pic.Session, job.Customer.ChatId)
		if err != nil {
			// http.Error(w, "Error sending typing event", http.StatusConflict)
			failedJob = append(failedJob, models.JobResponse{
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
			failedJob = append(failedJob, models.JobResponse{
				CustomerNumber: job.Customer.FormattedPhoneNumber,
				Name:           job.Customer.Name,
			})
			continue
		}
		time.Sleep(util.GenerateRandomDuration(30))

		//	send message
		msg := strings.ReplaceAll(TEXT_TEMPLATE, "{{name}}", job.Customer.Name)
		err = service.SendMessage(job.Pic.Session, job.Customer.ChatId, msg)
		if err != nil {
			failedJob = append(failedJob, models.JobResponse{
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(failedJob)
}
