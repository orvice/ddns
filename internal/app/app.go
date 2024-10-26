package app

import (
	"log/slog"

	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/config"
	"github.com/orvice/ddns/internal/ip"
)

type App struct {
	logger      *slog.Logger
	config      *config.Config
	dnsProvider dns.LibDNS
	ipGetter    ip.Getter
}

func New(logger *slog.Logger, config *config.Config, dnsProvider dns.LibDNS, ipGetter ip.Getter) *App {
	return &App{
		logger:      logger,
		config:      config,
		dnsProvider: dnsProvider,
		ipGetter:    ipGetter,
	}
}
