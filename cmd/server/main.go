package main

import (
	"log"
	"net/http"
	"waha-job-processing/internal/handler"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handler.Ping)
	mux.HandleFunc("POST /api/process-job", handler.ProcessJobHandler)
	mux.HandleFunc("GET /api/tracked-promo/{hash}", handler.GetTrackedPromos)
	mux.HandleFunc("POST /api/tracked-promo/{hash}/claim", handler.ClaimTrackedPromo)
	mux.HandleFunc("POST /api/log-blast", handler.CreateLogBlast)
	mux.HandleFunc("PATCH /api/log-blast", handler.UpdateLogBlast)

	log.Println("🚀 Server is running at http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// TODO: implement log blast, phone number not exist, and other features
// - log blast: track the status of each blast job
// - log blast created when starting the parent job, until all child jobs are finished, it should be updated again

// - phone number not exist: handle cases where the phone number is not found
// - other features: implement any additional features as needed

// note: n8n getting more complex. consider moving the preprocessing logic here
