package util

import (
	"errors"
	"log"
)

func ValidateWahaInputParams(session, chatId string) error {
	log.Printf("Validating input parameters: session=%s, chatId=%s", session, chatId)
	if session == "" || chatId == "" {
		return errors.New("missing session or chatId")
	}
	return nil
}
