package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"innogen-backend/run_service/internal/handler"
	"innogen-backend/run_service/internal/repository"
	"innogen-backend/run_service/internal/route"
	"innogen-backend/run_service/internal/service"
	"innogen-backend/shared/config"
	"innogen-backend/shared/database"
	"innogen-backend/shared/logger"
	"innogen-backend/shared/piston"
	"innogen-backend/shared/response"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()
	log := logger.New("run-service")

	ctx := context.Background()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	repo := repository.New(pool)
	pistonClient := piston.NewClient(cfg.PistonBaseURL)
	svc := service.New(repo, pistonClient)
	h := handler.New(svc, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	route.Register(mux, h)

	addr := ":" + cfg.RunServicePort
	log.Info("run-service listening on " + addr)

	handler := corsMiddleware(mux)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Error("server terminated with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// healthHandler responds with the service health status.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "run-service",
	})
}
