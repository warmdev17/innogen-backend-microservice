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
	log := logger.New("api-gateway")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(log))

	addr := ":" + cfg.APIGatewayPort
	log.Info("api-gateway listening on " + addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("server failed", slog.String("error", err.Error()))
	}
}

// healthHandler returns an http.HandlerFunc that responds with a JSON health
// check payload.
func healthHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("health check received")

		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "api-gateway",
		})
	}
}
