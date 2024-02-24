package main

import (
	"context"
	"fmt"
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))

	logger.Info("Hello, world!")

	return nil
}
