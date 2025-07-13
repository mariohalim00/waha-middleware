package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	"waha-job-processing/internal/models"
)

var secret = []byte("le-secret")

func BuildSignature(data models.PromoToken, timestamp string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jsonData)
	hashHex := hex.EncodeToString(hash[:])
	stringToSign := strings.ToLower(hashHex) + ":" + timestamp

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)

	return hex.EncodeToString(signature), nil
}

func VerifyingSignature(secret []byte, token models.PromoToken, timestamp, providedSignature string) error {
	expectedSignature, err := BuildSignature(token, timestamp)
	if err != nil {
		return err
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return errors.New("invalid timestamp")
	}

	if time.Now().Unix()-ts > 3600 {
		return errors.New("expired token")
	}

	// Constant-time compare
	if !hmac.Equal([]byte(expectedSignature), []byte(providedSignature)) {
		return errors.New("signature mismatch")
	}

	return nil
}
