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
	log := logger.New("repo-service")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "repo-service",
		})
	})

	addr := ":" + cfg.RepoServicePort
	log.Info("repo-service listening on " + addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("server stopped", slog.Any("error", err))
	}
}
