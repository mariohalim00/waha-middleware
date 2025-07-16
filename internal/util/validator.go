package util

import "errors"

func ValidateWahaInputParams(session, chatId string) error {
	if session == "" || chatId == "" {
		return errors.New("missing session or chatId")
	}
	return nil
}
