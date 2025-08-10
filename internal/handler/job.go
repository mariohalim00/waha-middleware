package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"waha-job-processing/internal/database/repository"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/service"
	"waha-job-processing/internal/service/waha"
	"waha-job-processing/internal/util"
	"waha-job-processing/internal/util/httpHelper"

	"github.com/joho/godotenv"
)

var TEXT_TEMPLATE string
var PROMO_CODE string
var WEBFORM_URL string
var MAX_CONCURRENT_JOBS int

func init() {
	godotenv.Load()
	TEXT_TEMPLATE = os.Getenv("TEXT_BLAST_TEMPLATE")
	PROMO_CODE = os.Getenv("PROMO_CODE")
	WEBFORM_URL = os.Getenv("BASE_WEBFORM_URL")

	// Set default concurrent jobs to 5, or read from environment
	MAX_CONCURRENT_JOBS = 5
	if envConcurrency := os.Getenv("MAX_CONCURRENT_JOBS"); envConcurrency != "" {
		if parsed, err := strconv.Atoi(envConcurrency); err == nil && parsed > 0 {
			MAX_CONCURRENT_JOBS = parsed
		}
	}
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
	// Configure concurrency - use environment variable or default
	maxWorkers, err := waha.GetAllActiveSessions()

	if err != nil {
		maxWorkers = min(len(jobList), MAX_CONCURRENT_JOBS)
	}
	log.Printf("Processing %d jobs with %d concurrent workers (workers here is session)", len(jobList), maxWorkers)

	// Channels for job distribution and result collection
	jobChan := make(chan models.Job, len(jobList))
	resultChan := make(chan *models.JobResponse, len(jobList))

	// Start worker goroutines
	for range maxWorkers {
		go jobWorker(jobChan, resultChan)
	}

	// Send jobs to workers
	for _, job := range jobList {
		jobChan <- job
	}
	close(jobChan)

	// Collect results
	var failedJobs []models.JobResponse
	for range jobList {
		result := <-resultChan
		if result != nil {
			failedJobs = append(failedJobs, *result)
		}
	}
	close(resultChan) // Close the result channel after collecting all results

	log.Printf("Job processing completed. Failed jobs: %d/%d", len(failedJobs), len(jobList))

	body, err := json.Marshal(failedJobs)
	if err != nil {
		log.Printf("Error converting failedJobs to JSON: %v. FailedJobs content: %+v\n", err, failedJobs)
	}

	err = callWebhookPostJob(body)
	if err != nil {
		log.Printf("Error when calling webhook: %+v\n", err)
	}
}

// jobWorker processes individual jobs from the job channel
func jobWorker(jobChan <-chan models.Job, resultChan chan<- *models.JobResponse) {
	for job := range jobChan {
		result := processIndividualJob(job)
		resultChan <- result
	}
}

// processIndividualJob handles the processing of a single job
func processIndividualJob(job models.Job) *models.JobResponse {
	startTime := time.Now()
	log.Println("[PROCESSING JOB - START] Processing job for:", job.Customer.FormattedPhoneNumber, "Session:", job.Pic.Session, "StartTime:", startTime.Format(time.RFC3339))

	// Helper function to create failed job response
	createFailedResponse := func() *models.JobResponse {
		return &models.JobResponse{
			CustomerNumber: job.Customer.FormattedPhoneNumber,
			Name:           job.Customer.Name,
			Voucher:        job.Customer.Voucher,
			Username:       job.Customer.Username,
		}
	}

	custdata, err := service.GetUserDataByUsername(job.Customer.Username)
	if err != nil {
		log.Println("Error fetching user data:", err)
		return createFailedResponse()
	}

	assignedVoucher, err := service.GetVoucherByName(job.Customer.Voucher)
	if err != nil {
		log.Printf("Error retrieving voucher: %v", err)
		return createFailedResponse()
	}

	// Start typing
	err = waha.StartTyping(job.Pic.Session, job.Customer.FormattedPhoneNumber)
	if err != nil {
		log.Printf("Error sending typing event: %v", err)
		return createFailedResponse()
	}
	time.Sleep(util.GenerateRandomDuration(15, 30))

	// Stop typing
	err = waha.StopTyping(job.Pic.Session, job.Customer.FormattedPhoneNumber)
	if err != nil {
		log.Printf("Error stopping typing event: %v", err)
		return createFailedResponse()
	}
	time.Sleep(util.GenerateRandomDuration(15, 30))

	// Generate URL
	url, err := generateWebFormUrl(job, assignedVoucher, custdata.Userid)
	if err != nil {
		log.Printf("failed to generate url, skipping")
		return createFailedResponse()
	}

	// Prepare message template
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

	// Send message
	err = waha.SendMessage(job.Pic.Session, job.Customer.FormattedPhoneNumber, msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return createFailedResponse()
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Println("[PROCESSING JOB - FINISH] Processing job for:", job.Customer.FormattedPhoneNumber, "Session:", job.Pic.Session, "EndTime:", endTime.Format(time.RFC3339), "Duration:", duration)

	// Add random delay to avoid overwhelming the service
	time.Sleep(util.GenerateRandomDuration(10, 20))
	// Return nil for successful jobs (no failed response needed)
	return nil
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
