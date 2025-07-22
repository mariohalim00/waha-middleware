package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"waha-job-processing/internal/database/db"
	"waha-job-processing/internal/database/repository"
	"waha-job-processing/internal/util/httpHelper"
)

func GetTrackedPromos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	dbConn := db.New(ctx)
	q := repository.New(dbConn)

	sluggedHash := r.PathValue("hash")

	log.Println("Fetching tracked promo for hash:", sluggedHash)
	promos, err := q.GetOneTrackedPromo(ctx, sluggedHash)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to retrieve tracked promos", http.StatusInternalServerError)
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

func ClaimTrackedPromo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	dbConn := db.New(ctx)
	q := repository.New(dbConn)

	sluggedHash := r.PathValue("hash")

	existingPromo, err := q.GetOneTrackedPromo(ctx, sluggedHash)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to retrieve tracked promo", http.StatusInternalServerError)
		return
	}

	if existingPromo.Claimed {
		httpHelper.ReturnHttpError(w, "Promo already claimed", http.StatusConflict)
		return
	}

	log.Println("Claiming tracked promo for hash:", sluggedHash)

	param := repository.UpdateTrackedPromoParams{
		HashedString: sluggedHash,
		Claimed:      true,
	}
	res, err := q.UpdateTrackedPromo(ctx, param)
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

func MarkPromoAsProcessed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpHelper.ReturnHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	dbConn := db.New(ctx)
	defer dbConn.Close(ctx)

	q := repository.New(dbConn)
	sluggedHash := r.PathValue("hash")
	log.Println("Marking promo as processed for hash:", sluggedHash)
	param := repository.UpdatePromoTrackerIsProcessedParams{
		HashedString: sluggedHash,
		IsProcessed:  true,
	}
	_, err := q.UpdatePromoTrackerIsProcessed(ctx, param)
	if err != nil {
		httpHelper.ReturnHttpError(w, "Failed to mark promo as processed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
