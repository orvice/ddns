// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/app"
	"github.com/orvice/ddns/internal/config"
	"github.com/orvice/ddns/internal/ip"
	"github.com/orvice/ddns/notify"
	"github.com/orvice/ddns/utils"
)

// Injectors from wire.go:

func NewApp() (*app.App, error) {
	logger := utils.NewLogger()
	configConfig, err := config.New()
	if err != nil {
		return nil, err
	}
	libDNS := dns.New(configConfig)
	getter := ip.NewIfGetter()
	notifier, err := notify.NewTelegramNotifier(configConfig)
	if err != nil {
		return nil, err
	}
	appApp := app.New(logger, configConfig, libDNS, getter, notifier)
	return appApp, nil
}
