package service

import (
	"log"
	"os"
	"waha-job-processing/internal/util/httpHelper"
)

var TM_API_BASE_URL = os.Getenv("TM_API_URL")

func UpdateLastBlastByUsername(userName string) error {
	url := TM_API_BASE_URL + "/api/custdata/update-last-blast-date/" + userName
	err := httpHelper.Post([]byte{}, url, httpHelper.HttpHeader{})

	if err != nil {
		log.Println("Error executing request:", err)
		return err
	}

	return nil
}
