package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"innogen-backend/shared/config"
	"innogen-backend/shared/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New("submission-worker")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Info("submission-worker started", "redisAddr", cfg.RedisAddr)

	<-ctx.Done()

	log.Info("submission-worker stopped")
}
