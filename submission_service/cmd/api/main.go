package main

import (
	"innogen-backend/shared/config"
	"innogen-backend/shared/logger"
	"innogen-backend/shared/response"
	"log/slog"
	"net/http"
)

func main() {
	cfg := config.Load()
	log := logger.New("submission-service")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(log))

	addr := ":" + cfg.SubmissionServicePort
	log.Info("submission-service listening", slog.String("addr", addr))

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("server terminated with error", slog.String("error", err.Error()))
	}
}

func healthHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "submission-service",
		})
	}
}
