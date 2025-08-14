package handler

import (
	"fmt"
	"net/http"
	"waha-job-processing/internal/service"
	logblast "waha-job-processing/internal/service/log-blast"
	trackedpromo "waha-job-processing/internal/service/tracked-promo"
	"waha-job-processing/internal/service/vouchers"
)

type Handler struct {
	TrackedPromoService *trackedpromo.TrackedPromoService
	LogBlastService     *logblast.LogBlastService
	VoucherService      *vouchers.VoucherService
}

func NewHandler(services *service.Services) *Handler {
	handlerWithService := &Handler{
		TrackedPromoService: services.TrackedPromo,
		LogBlastService:     services.LogBlast,
		VoucherService:      services.Voucher,
	}

	return handlerWithService
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "pong"}`)
}
