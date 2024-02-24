package main

import (
	"context"
	"github.com/lmittmann/tint"
	"godns/dns"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))

	srv := dns.NewServer(logger, ctx)
	addr, err := srv.Start()

	logger.Info("server started", "addr", addr)

	if err != nil {
		logger.Error("error on serve", "err", err)
	}

	logger.Info("server stopped")
	return nil
}
