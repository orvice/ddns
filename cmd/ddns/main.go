package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/orvice/ddns/internal/wire"
)

func main() {
	app, err := wire.NewApp()
	if err != nil {
		slog.Error("init app error", "error", err)
		os.Exit(1)
	}
	if err := app.Run(context.Background()); err != nil {
		slog.Error("app run error", "error", err)
		os.Exit(1)
	}
}
