package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	port int
	db   *pgxpool.Pool
}

func NewServer(dbConn *pgxpool.Pool) *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Printf("Error parsing PORT environment variable: %v. Using default port 8080.", err)
		port = 8080 // Default port if not set in environment
	}

	NewServer := &Server{
		port: port,
		db:   dbConn,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(dbConn),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
