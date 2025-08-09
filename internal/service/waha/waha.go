package waha

import (
	"encoding/json"
	"log"
	"os"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/util"
	"waha-job-processing/internal/util/httpHelper"
)

var WAHA_BASE_URL = os.Getenv("WAHA_URL")
var WEBFORM_URL = os.Getenv("BASE_WEBFORM_URL")
var WAHA_HTTP_HEADER_POST = httpHelper.HttpHeader{
	"Content-Type": "application/json",
	"X-Api-Key":    os.Getenv("API_KEY"),
}

func StartTyping(session, chatId string) error {
	err := util.ValidateWahaInputParams(session, chatId)

	if err != nil {
		log.Println("Error: start typing validation", err)
		return err
	}

	chatDetails := models.ChatDetails{
		Session: session,
		ChatId:  chatId,
	}

	payload, err := json.Marshal(chatDetails)
	if err != nil {
		log.Println("Error trying to convert payload to JSON")
		return err
	}

	url := WAHA_BASE_URL + "/api/startTyping"

	err = httpHelper.Post(payload, url, WAHA_HTTP_HEADER_POST)

	if err != nil {
		log.Println("Error executing request:", err)
		return err
	}

	return nil
}

func StopTyping(session, chatId string) error {
	err := util.ValidateWahaInputParams(session, chatId)

	if err != nil {
		log.Println("Error: start typing validation", err)
		return err
	}

	chatDetails := models.ChatDetails{
		Session: session,
		ChatId:  chatId,
	}

	payload, err := json.Marshal(chatDetails)
	if err != nil {
		log.Println("Error trying to convert payload to JSON")
		return err
	}

	url := WAHA_BASE_URL + "/api/stopTyping"

	err = httpHelper.Post(payload, url, WAHA_HTTP_HEADER_POST)

	if err != nil {
		log.Println("Error executing request:", err)
		return err
	}

	return nil
}

func SendMessage(session, chatId, text string) error {
	err := util.ValidateWahaInputParams(session, chatId)

	if err != nil {
		log.Println("Error: start typing validation", err)
		return err
	}

	chatDetails := models.SendTextDetails{
		Session:                session,
		ChatId:                 chatId,
		ReplyTo:                "",
		LinkPreview:            true,
		LinkPreviewHighQuality: false,
		Text:                   text,
	}

	payload, err := json.Marshal(chatDetails)
	if err != nil {
		log.Println("Error trying to convert payload to JSON")
		return err
	}

	url := WAHA_BASE_URL + "/api/sendText"

	err = httpHelper.Post(payload, url, WAHA_HTTP_HEADER_POST)

	if err != nil {
		log.Println("Error executing request:", err)

	}

	return nil
}

func GetAllActiveSessions() (int, error) {
	url := WAHA_BASE_URL + "/api/sessions"
	resp, err := httpHelper.Get(url, WAHA_HTTP_HEADER_POST)
	if err != nil {
		log.Println("Error fetching active sessions:", err)
		return -1, err
	}

	var sessions []WahaSession
	if err := json.Unmarshal(resp, &sessions); err != nil {
		log.Println("Error decoding response:", err)
		return -1, err
	}

	workingCount := 0

	for _, session := range sessions {
		if session.Status == StatusWorking {
			workingCount++
		}
	}

	log.Println("Active sessions fetched successfully:", workingCount)
	return workingCount, nil
}
