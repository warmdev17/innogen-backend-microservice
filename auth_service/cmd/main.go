package main

import (
	"log/slog"
	"net/http"

	"innogen-backend/shared/config"
	"innogen-backend/shared/logger"
	"innogen-backend/shared/response"
)

func main() {
	cfg := config.Load()
	log := logger.New("auth-service")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	addr := ":" + cfg.AuthServicePort
	log.Info("auth-service listening on " + addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("server terminated with error", slog.String("error", err.Error()))
	}
}

// healthHandler responds with the service health status.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "auth-service",
	})
}
