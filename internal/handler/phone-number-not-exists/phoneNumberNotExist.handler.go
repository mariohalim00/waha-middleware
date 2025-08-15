package phonenumbernotexists

import (
	"encoding/json"
	"net/http"
	"waha-job-processing/internal/handler"
	"waha-job-processing/internal/util/httpHelper"
)

type Handler struct {
	*handler.Handler
}
type CreatePhoneNumberNotExistDto struct {
	PhoneNumber string `json:"phoneNumber"`
	Username    string `json:"username"`
	BlastID     string `json:"blastId"`
}

func (h *Handler) CreatePhoneNumberNotExist(w http.ResponseWriter, r *http.Request) {
	var phoneNumberDto CreatePhoneNumberNotExistDto
	err := json.NewDecoder(r.Body).Decode(&phoneNumberDto)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if phoneNumberDto.PhoneNumber == "" || phoneNumberDto.Username == "" || phoneNumberDto.BlastID == "" {
		httpHelper.ReturnHttpError(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	res, err := h.PhoneNumberNotExistService.CreatePhoneNumberNotExist(phoneNumberDto.PhoneNumber, phoneNumberDto.Username, phoneNumberDto.BlastID)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to create phone number not exist entry", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
