package util

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

var API_KEY = os.Getenv("API_KEY")

func Post(payload []byte, url string) error {

	//TODO: REMOVE DEBUG
	fmt.Println("[http/POST] url: ", url)
	fmt.Println("[http/POST] payload: ", payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", API_KEY)

	client := &http.Client{}

	_, err = client.Do(req)

	if err != nil {
		fmt.Println("Error processing HTTP Request", err)
		return err
	}

	return nil
}
