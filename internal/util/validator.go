package util

import "errors"

func ValidateWahaInputParams(session, chatId string) error {
	if len(session) == 0 || len(chatId) == 0 {
		return errors.New("Missing session or chatId")
	}

	return nil
}
