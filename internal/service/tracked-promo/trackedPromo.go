package trackedpromo

import (
	"context"
	"errors"
	"log"
	"time"
	"waha-job-processing/internal/database/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackedPromoService struct {
	DB *pgxpool.Pool
}

func NewTrackedPromoService(database *pgxpool.Pool) *TrackedPromoService {
	return &TrackedPromoService{DB: database}
}

func (s *TrackedPromoService) CreateTrackedPromo(signature, userName, voucher string, expiryDate time.Time, userID string) (repository.PromoTracker, error) {
	ctx := context.Background()
	query := repository.New(s.DB)

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

	trackedPromo, err := query.CreateTrackedPromo(ctx, param)
	if err != nil {
		log.Printf("Failed to crate tracked promo record")
		return repository.PromoTracker{}, err
	}

	log.Printf("Created tracked promo %v\n", trackedPromo)
	return trackedPromo, nil
}

func (s *TrackedPromoService) GetTrackedPromo(hash string) (repository.PromoTracker, error) {
	ctx := context.Background()
	q := repository.New(s.DB)

	log.Println("Fetching tracked promo for hash:", hash)
	promos, err := q.GetOneTrackedPromo(ctx, hash)
	if err != nil {
		log.Printf("Failed to retrieve tracked promo for hash %s: %v", hash, err)
		return repository.PromoTracker{}, err
	}

	return promos, nil
}

func (s *TrackedPromoService) ClaimTrackedPromo(hash string) (repository.PromoTracker, error) {
	ctx := context.Background()
	q := repository.New(s.DB)

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

func (s *TrackedPromoService) MarkPromoAsProcessed(hash string) (repository.PromoTracker, error) {
	ctx := context.Background()
	q := repository.New(s.DB)

	existingPromo, err := s.GetTrackedPromo(hash)
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
