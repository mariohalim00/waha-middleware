package trackedpromo

import (
	"encoding/json"
	"net/http"
	"waha-job-processing/internal/handler"
	"waha-job-processing/internal/util/httpHelper"
)

// Handler is a local type that embeds handler.Handler.
type Handler struct {
	*handler.Handler
}

func (h *Handler) GetTrackedPromos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sluggedHash := r.PathValue("hash")

	promos, err := h.TrackedPromoService.ClaimTrackedPromo(sluggedHash)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to retrieve tracked promo", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(promos)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) ClaimTrackedPromo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sluggedHash := r.PathValue("hash")

	res, err := h.TrackedPromoService.ClaimTrackedPromo(sluggedHash)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to claim tracked promo", http.StatusInternalServerError)
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

func (h *Handler) MarkPromoAsProcessed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sluggedHash := r.PathValue("hash")

	_, err := h.TrackedPromoService.MarkPromoAsProcessed(sluggedHash)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to mark promo as processed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Promo marked as processed successfully",
	})
}
