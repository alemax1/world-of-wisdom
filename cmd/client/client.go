package main

import (
	"flag"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/alemax1/world-of-wisdom/internal/client"
	"github.com/alemax1/world-of-wisdom/pkg/pow"
)

func main() {
	loggerConfig := zap.NewDevelopmentEncoderConfig()
	loggerConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(loggerConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	))

	var addr string

	flag.StringVar(&addr, "addr", "localhost:9999", "connect address")
	flag.Parse()

	solver := pow.NewSolver()

	cl := client.New(solver, addr, logger)

	if err := cl.Run(); err != nil {
		logger.Error("run client", zap.Error(err))
	}
}
