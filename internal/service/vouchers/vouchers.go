package vouchers

import (
	"context"
	"log"
	"waha-job-processing/internal/database/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoucherService struct {
	DB *pgxpool.Pool
}

func NewVoucherService(database *pgxpool.Pool) *VoucherService {
	return &VoucherService{DB: database}
}

func (s *VoucherService) GetVoucherByName(name string) (*repository.Voucher, error) {
	ctx := context.Background()
	query := repository.New(s.DB)

	voucher, err := query.GetOneVoucher(ctx, name)
	if err != nil {
		log.Printf("Failed to get voucher by name: %v", name)
		return nil, err
	}

	log.Printf("Retrieved voucher: %v", voucher)
	return &voucher, nil
}
