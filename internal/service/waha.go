package service

import (
	"encoding/json"
	"log"
	"os"
	"waha-job-processing/internal/models"
	"waha-job-processing/internal/util"
	"waha-job-processing/internal/util/httpHelper"
)

var BASE_URL = os.Getenv("WAHA_URL")
var WEBFORM_URL = os.Getenv("BASE_WEBFORM_URL")

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

	url := BASE_URL + "/api/startTyping"

	err = httpHelper.Post(payload, url)

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

	url := BASE_URL + "/api/stopTyping"

	err = httpHelper.Post(payload, url)

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

	url := BASE_URL + "/api/sendText"

	err = httpHelper.Post(payload, url)

	if err != nil {
		log.Println("Error executing request:", err)

	}

	return nil
}
