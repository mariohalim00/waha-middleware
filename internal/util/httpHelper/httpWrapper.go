package httpHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HttpHeader map[string]string

// Helper function to check for successful status codes
func isSuccessfulStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

// Post sends a POST request with the provided payload to the specified URL.
func Post(payload []byte, url string, headers HttpHeader) error {
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return err
	}

	// Set headers from parameters
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// HTTP Client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // 10 seconds timeout
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error processing HTTP request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check for successful response status code
	if !isSuccessfulStatusCode(resp.StatusCode) {
		// Read and log the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read error response body: %v", err)
			return fmt.Errorf("HTTP error: %s", resp.Status)
		}
		log.Printf("Error Response: %s", string(body))
		return fmt.Errorf("HTTP error: %s\nBody: %s", resp.Status, string(body))
	}

	return nil
}

func Get(url string, headers HttpHeader) ([]byte, error) {
	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return nil, err
	}

	// Set headers from parameters
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// HTTP Client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // 10 seconds timeout
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error processing HTTP request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for successful response status code
	if !isSuccessfulStatusCode(resp.StatusCode) {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error: %s\nBody: %s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}

type HttpError struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func ReturnHttpError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	response := HttpError{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Message:    message,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
