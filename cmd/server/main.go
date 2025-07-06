package main

import (
	"fmt"
	"log"
	"net/http"
	"waha-job-processing/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handler.Ping)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("🚀 Server is running at http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
