package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Post(payload []byte, url string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", os.Getenv("API_KEY"))

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error processing HTTP Request", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error: %s\nBody: %s", resp.Status, string(body))
	}

	return nil
}
