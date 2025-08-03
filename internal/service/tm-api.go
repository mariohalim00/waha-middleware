package service

import (
	"encoding/json"
	"log"
	"os"
	"waha-job-processing/internal/models"
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

func GetUserDataByUsername(username string) (models.CustdatumDto, error) {
	url := TM_API_BASE_URL + "/api/custdata/" + username
	response, err := httpHelper.Get(url, httpHelper.HttpHeader{})

	if err != nil {
		log.Println("Error executing request:", err)
		return models.CustdatumDto{}, err
	}

	var userData models.CustdatumDto
	err = json.Unmarshal(response, &userData)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return models.CustdatumDto{}, err
	}

	return userData, nil
}
