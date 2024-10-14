package main

import (
	"flag"

	"go.uber.org/zap"

	"github.com/alemax1/world-of-wisdom/internal/client"
	"github.com/alemax1/world-of-wisdom/pkg/logger"
	"github.com/alemax1/world-of-wisdom/pkg/pow"
)

func main() {
	log := logger.New()

	var addr string

	flag.StringVar(&addr, "addr", "localhost:9999", "connect address")
	flag.Parse()

	solver := pow.NewSolver()

	cl := client.New(solver, addr, log)

	if err := cl.Run(); err != nil {
		log.Error("run client", zap.Error(err))
	}
}
