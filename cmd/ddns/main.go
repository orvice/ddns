package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/ip"
	"github.com/orvice/ddns/internal/wire"
)

var (
	dnsProvider dns.LibDNS
	ipGetter    ip.Getter
	DNSMode     string
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
