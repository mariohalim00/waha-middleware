package service

import (
	"context"
	"log"
	"waha-job-processing/internal/database/db"
	"waha-job-processing/internal/database/repository"
)

func GetVoucherByName(name string) (*repository.Voucher, error) {
	ctx := context.Background()
	conn := db.New(ctx)
	defer conn.Close(ctx)

	query := repository.New(conn)

	voucher, err := query.GetOneVoucher(ctx, name)
	if err != nil {
		log.Printf("Failed to get voucher by name: %v", name)
		return nil, err
	}

	log.Printf("Retrieved voucher: %v", voucher)
	return &voucher, nil
}
