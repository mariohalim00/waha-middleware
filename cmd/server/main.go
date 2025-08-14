package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"waha-job-processing/internal/database/db"
	"waha-job-processing/internal/server"

	"github.com/joho/godotenv"
)

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	ctx := context.Background()
	dbConnection := db.New(ctx)
	defer dbConnection.Close()

	// Initializer server
	server := server.NewServer(dbConnection)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	log.Printf("[SERVER INITIALIZED] 🚀 Server is running at http://localhost:%v", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}

	<-done
	log.Println("Server shutdown gracefully")

}

// TODO: implement log blast, phone number not exist, and other features
// - log blast: track the status of each blast job
// - log blast created when starting the parent job, until all child jobs are finished, it should be updated again

// - phone number not exist: handle cases where the phone number is not found
// - other features: implement any additional features as needed

// note: n8n getting more complex. consider moving the preprocessing logic here
