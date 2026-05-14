package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"innogen-backend/shared/config"
	"innogen-backend/shared/logger"
	"innogen-backend/shared/response"
)

func main() {
	cfg := config.Load()
	log := logger.New("run-service")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "run-service",
		})
	})

	addr := fmt.Sprintf(":%s", cfg.RunServicePort)
	log.Info("run-service listening on " + addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("server failed", slog.Any("error", err))
	}
}
