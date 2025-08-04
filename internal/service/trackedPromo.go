package service

import (
	"context"
	"errors"
	"log"
	"time"
	"waha-job-processing/internal/database/db"
	"waha-job-processing/internal/database/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

func CreateTrackedPromo(signature, userName, voucher string, expiryDate time.Time, userID string) (repository.PromoTracker, error) {
	ctx := context.Background()
	conn := db.New(ctx)
	defer conn.Close(ctx)

	param := repository.CreateTrackedPromoParams{
		HashedString: signature,
		ExpiredAt: pgtype.Timestamptz{
			Time:  expiryDate,
			Valid: true,
		},
		UserName: userName,
		Voucher: pgtype.Text{
			String: voucher,
			Valid:  true,
		},
		UserID: pgtype.Text{
			String: userID,
			Valid:  true,
		},
	}

	query := repository.New(conn)
	trackedPromo, err := query.CreateTrackedPromo(ctx, param)
	if err != nil {
		log.Printf("Failed to crate tracked promo record")
		return repository.PromoTracker{}, err
	}

	log.Printf("Created tracked promo %v\n", trackedPromo)
	return trackedPromo, nil
}

func GetTrackedPromo(hash string) (repository.PromoTracker, error) {
	ctx := context.Background()
	dbConn := db.New(ctx)
	defer dbConn.Close(ctx)
	q := repository.New(dbConn)

	log.Println("Fetching tracked promo for hash:", hash)
	promos, err := q.GetOneTrackedPromo(ctx, hash)
	if err != nil {
		log.Printf("Failed to retrieve tracked promo for hash %s: %v", hash, err)
		return repository.PromoTracker{}, err
	}

	return promos, nil
}

func ClaimTrackedPromo(hash string) (repository.PromoTracker, error) {
	ctx := context.Background()
	dbConn := db.New(ctx)
	defer dbConn.Close(ctx)
	q := repository.New(dbConn)

	existingPromo, err := q.GetOneTrackedPromo(ctx, hash)
	if err != nil {
		log.Printf("Failed to retrieve tracked promo for hash %s: %v", hash, err)
		return repository.PromoTracker{}, err
	}

	if existingPromo.Claimed {
		return repository.PromoTracker{}, err
	}

	log.Println("Claiming tracked promo for hash:", hash)

	param := repository.UpdateTrackedPromoParams{
		HashedString: hash,
		Claimed:      true,
	}
	res, err := q.UpdateTrackedPromo(ctx, param)
	if err != nil {
		log.Printf("Failed to claim tracked promo for hash %s: %v", hash, err)
		return repository.PromoTracker{}, err
	}

	return res, nil
}

func MarkPromoAsProcessed(hash string) (repository.PromoTracker, error) {
	ctx := context.Background()
	dbConn := db.New(ctx)
	defer dbConn.Close(ctx)

	q := repository.New(dbConn)

	existingPromo, err := GetTrackedPromo(hash)
	if err != nil {
		log.Printf("Failed to retrieve tracked promo for hash %s: %v", hash, err)
		return repository.PromoTracker{}, err
	}

	if existingPromo.IsProcessed {
		log.Printf("Promo with hash %s is already processed", hash)
		err = errors.New("promo already processed")
		return existingPromo, err
	}

	log.Println("Marking promo as processed for hash:", hash)
	param := repository.UpdatePromoTrackerIsProcessedParams{
		HashedString: hash,
		IsProcessed:  true,
	}

	res, err := q.UpdatePromoTrackerIsProcessed(ctx, param)
	if err != nil {
		log.Printf("[Service] Failed to mark promo as processed for hash %s: %v", hash, err)
		return repository.PromoTracker{}, err
	}

	return res, nil
}
