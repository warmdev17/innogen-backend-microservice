package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	reporepo "innogen-backend/repo_service/repository"
	reposervice "innogen-backend/repo_service/service"
	"innogen-backend/shared/config"
	"innogen-backend/shared/database"
	"innogen-backend/shared/logger"
	"innogen-backend/shared/piston"
	"innogen-backend/submission_service/internal/queue"
	"innogen-backend/submission_service/internal/repository"
	"innogen-backend/submission_service/internal/worker"
)

func main() {
	cfg := config.Load()
	log := logger.New("submission-worker")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	q, err := queue.New(cfg.RedisAddr)
	if err != nil {
		log.Error("failed to connect to redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer q.Close()

	repo := repository.New(pool)
	repoRepo := reporepo.New(pool)
	repoSvc := reposervice.New(repoRepo, log)
	pistonClient := piston.NewClient(cfg.PistonBaseURL)
	w := worker.New(log, repo, pistonClient, q, repoSvc)

	log.Info("submission-worker started")
	w.Run(ctx)
	log.Info("submission-worker stopped")
}
