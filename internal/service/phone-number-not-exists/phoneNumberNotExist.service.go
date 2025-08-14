package phonenumbernotexists

import (
	"context"
	"log"
	"waha-job-processing/internal/database/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhoneNumberNotExistService struct {
	DB *pgxpool.Pool
}

func NewPhoneNumberNotExistService(database *pgxpool.Pool) *PhoneNumberNotExistService {
	return &PhoneNumberNotExistService{DB: database}
}

func (s *PhoneNumberNotExistService) CreatePhoneNumberNotExist(phoneNumber, username string, blastID string) (repository.PhoneNumberNotExist, error) {
	query := repository.New(s.DB)

	blastUuid := pgtype.UUID{}
	err := blastUuid.Scan(blastID)

	ctx := context.Background()

	if err != nil {
		log.Printf("Failed to scan blast ID: %v", err)
		return repository.PhoneNumberNotExist{}, err
	}

	phoneNumberNotExist := repository.CreatePhoneNumberNotExistParams{
		PhoneNumber: phoneNumber,
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
		BlastID: blastUuid,
	}

	res, err := query.CreatePhoneNumberNotExist(ctx, phoneNumberNotExist)
	if err != nil {
		return repository.PhoneNumberNotExist{}, err
	}

	return res, nil
}
