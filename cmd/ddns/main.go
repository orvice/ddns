package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/orvice/ddns/internal/wire"
)

var (
	IPNotifyFormat = "[%s] ip changed, old IP: %s new IP: %s"
)

func main() {
	app, err := wire.NewApp()
	if err != nil {
		slog.Error("init app error", "error", err)
		os.Exit(1)
	}
	app.Run(context.Background())
}
