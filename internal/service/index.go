package service

import (
	logblast "waha-job-processing/internal/service/log-blast"
	trackedpromo "waha-job-processing/internal/service/tracked-promo"
	"waha-job-processing/internal/service/vouchers"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	TrackedPromo *trackedpromo.TrackedPromoService
	LogBlast     *logblast.LogBlastService
	Voucher      *vouchers.VoucherService
}

func InitializeServices(database *pgxpool.Pool) *Services {
	trackedPromoService := trackedpromo.NewTrackedPromoService(database)
	logBlastService := logblast.NewLogBlastService(database)
	voucherService := vouchers.NewVoucherService(database)

	return &Services{
		TrackedPromo: trackedPromoService,
		LogBlast:     logBlastService,
		Voucher:      voucherService,
	}
}
