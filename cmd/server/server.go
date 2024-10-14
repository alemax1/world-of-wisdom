package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alemax1/world-of-wisdom/internal/config"
	"github.com/alemax1/world-of-wisdom/internal/server"
	"github.com/alemax1/world-of-wisdom/internal/storage"
	"github.com/alemax1/world-of-wisdom/pkg/logger"
	"github.com/alemax1/world-of-wisdom/pkg/pow"

	"go.uber.org/zap"
)

func main() {
	log := logger.New()
	cfg := config.New()

	challenger := pow.NewManager(cfg.ChallengeSize)

	st, err := storage.New(cfg.DataFilePath)
	if err != nil {
		log.Fatal("cannot create storage", zap.Error(err))
	}

	srv := server.New(challenger, st, *cfg, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		cancel()
	}()

	if err := srv.Run(ctx); err != nil {
		log.Error("server run", zap.Error(err))
	}
}
