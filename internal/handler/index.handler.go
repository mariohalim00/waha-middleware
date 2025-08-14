package handler

import (
	"fmt"
	"net/http"
	"waha-job-processing/internal/service"
	logblast "waha-job-processing/internal/service/log-blast"
	phonenumbernotexists "waha-job-processing/internal/service/phone-number-not-exists"
	trackedpromo "waha-job-processing/internal/service/tracked-promo"
	"waha-job-processing/internal/service/vouchers"
)

type Handler struct {
	TrackedPromoService        *trackedpromo.TrackedPromoService
	LogBlastService            *logblast.LogBlastService
	VoucherService             *vouchers.VoucherService
	PhoneNumberNotExistService *phonenumbernotexists.PhoneNumberNotExistService
}

func NewHandler(services *service.Services) *Handler {
	handlerWithService := &Handler{
		TrackedPromoService:        services.TrackedPromo,
		LogBlastService:            services.LogBlast,
		VoucherService:             services.Voucher,
		PhoneNumberNotExistService: services.PhoneNumberNotExist,
	}

	return handlerWithService
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "pong"}`)
}
