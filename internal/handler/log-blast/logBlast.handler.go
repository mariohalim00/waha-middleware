package logblast

import (
	"encoding/json"
	"net/http"
	"time"
	"waha-job-processing/internal/database/repository"
	"waha-job-processing/internal/handler"
	"waha-job-processing/internal/util/httpHelper"

	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	*handler.Handler
}

type CreateLogBlastParamsDto struct {
	WorkflowStart     *time.Time `json:"workflow_start"`
	BlastStart        *time.Time `json:"blast_start"`
	BlastEnd          *time.Time `json:"blast_end"`
	ActualBlast       *int32     `json:"actual_blast"`
	SuccessBlast      *int32     `json:"success_blast"`
	FailedBlast       *int32     `json:"failed_blast"`
	RawBlast          *int32     `json:"raw_blast"`
	NonExistentNumber *int32     `json:"non_existent_number"`
}

func (h *Handler) CreateLogBlast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logBlastParamDto CreateLogBlastParamsDto
	err := json.NewDecoder(r.Body).Decode(&logBlastParamDto)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields before using them
	if logBlastParamDto.RawBlast == nil {
		httpHelper.ReturnHttpError(w, "Missing required fields: raw_blast", http.StatusBadRequest)
		return
	}

	createLogBlastParams := repository.CreateLogBlastParams{
		WorkflowStart: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		RawBlast: pgtype.Int4{
			Int32: *logBlastParamDto.RawBlast,
			Valid: true,
		},
		NonExistentNumber: pgtype.Int4{
			Int32: *logBlastParamDto.NonExistentNumber,
			Valid: true,
		},
	}
	res, err := h.LogBlastService.CreateLogBlast(createLogBlastParams)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to create log blast", http.StatusInternalServerError)
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

type UpdateLogBlastParamsDto struct {
	Id                string     `json:"id"`
	WorkflowStart     *time.Time `json:"workflow_start"`
	BlastStart        *time.Time `json:"blast_start"`
	BlastEnd          *time.Time `json:"blast_end"`
	ActualBlast       *int32     `json:"actual_blast"`
	SuccessBlast      *int32     `json:"success_blast"`
	FailedBlast       *int32     `json:"failed_blast"`
	RawBlast          *int32     `json:"raw_blast"`
	NonExistentNumber *int32     `json:"non_existent_number"`
}

func (h *Handler) UpdateLogBlast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logBlastParamDto UpdateLogBlastParamsDto
	err := json.NewDecoder(r.Body).Decode(&logBlastParamDto)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if logBlastParamDto.Id == "" {
		httpHelper.ReturnHttpError(w, "Missing required field: id", http.StatusBadRequest)
		return
	}

	updateUUID := pgtype.UUID{}
	err = updateUUID.Scan(logBlastParamDto.Id)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Invalid UUID format for id", http.StatusBadRequest)
		return
	}

	updateLogBlastParam := repository.UpdateLogBlastParams{
		ID: updateUUID,
		BlastStart: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		BlastEnd: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ActualBlast: pgtype.Int4{
			Int32: *logBlastParamDto.ActualBlast,
			Valid: true,
		},
		SuccessBlast: pgtype.Int4{
			Int32: *logBlastParamDto.SuccessBlast,
			Valid: true,
		},
		FailedBlast: pgtype.Int4{
			Int32: *logBlastParamDto.FailedBlast,
			Valid: true,
		},
	}

	res, err := h.LogBlastService.UpdateLogBlast(updateLogBlastParam)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to update log blast", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
