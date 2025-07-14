package main

import (
	"fmt"
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

	fmt.Println("🚀 Server is running at http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
