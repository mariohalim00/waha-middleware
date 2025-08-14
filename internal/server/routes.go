package server

import (
	"net/http"
	"waha-job-processing/internal/handler"
	"waha-job-processing/internal/handler/job"
	logblast "waha-job-processing/internal/handler/log-blast"
	phonenumbernotexists "waha-job-processing/internal/handler/phone-number-not-exists"
	trackedpromo "waha-job-processing/internal/handler/tracked-promo"
	"waha-job-processing/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Server) RegisterRoutes(dbConn *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	services := service.InitializeServices(dbConn)

	// top-level handler
	handlerWithService := handler.NewHandler(services)

	// register routes here
	mux.HandleFunc("GET /ping", handlerWithService.Ping)

	// Mount private routes under a prefix, e.g., "/api/"
	privateApiRoutes := s.registerPrivateRoutes(handlerWithService)
	mux.Handle("/api/", http.StripPrefix("/api", privateApiRoutes))

	return mux
}

func (s *Server) privateRouteMiddlewareWrapper(h http.Handler) http.Handler {
	return ChainMiddleware(h, s.CorsMiddleware, s.LogMiddleware)
}

func (s *Server) registerPrivateRoutes(mainHandler *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	// specific handlers goes here
	logBlastHandlers := &logblast.Handler{
		Handler: mainHandler,
	}

	trackedPromoHandlers := &trackedpromo.Handler{
		Handler: mainHandler,
	}

	jobHandler := &job.Handler{
		Handler: mainHandler,
	}

	phoneNumberNotExistHandlers := &phonenumbernotexists.Handler{
		Handler: mainHandler,
	}

	mux.HandleFunc("POST /process-job", jobHandler.ProcessJobHandler)

	mux.HandleFunc("GET /tracked-promo/{hash}", trackedPromoHandlers.GetTrackedPromos)
	mux.HandleFunc("POST /tracked-promo/{hash}/claim", trackedPromoHandlers.ClaimTrackedPromo)

	mux.HandleFunc("POST /log-blast", logBlastHandlers.CreateLogBlast)
	mux.HandleFunc("PATCH /log-blast", logBlastHandlers.UpdateLogBlast)

	mux.HandleFunc("GET /health", mainHandler.Ping)

	mux.HandleFunc("POST /phone-number-not-exists", phoneNumberNotExistHandlers.CreatePhoneNumberNotExist)

	return s.privateRouteMiddlewareWrapper(mux)
}
