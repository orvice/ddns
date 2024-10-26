package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/config"
	"github.com/orvice/ddns/internal/ip"
	"github.com/orvice/ddns/notify"
)

var (
	IPNotifyFormat = "[%s] ip changed, old IP: %s new IP: %s"
)

type App struct {
	logger      *slog.Logger
	config      *config.Config
	dnsProvider dns.LibDNS
	ipGetter    ip.Getter
	notifier    notify.Notifier
}

func New(logger *slog.Logger, config *config.Config, dnsProvider dns.LibDNS, ipGetter ip.Getter, notifier notify.Notifier) *App {
	return &App{
		logger:      logger,
		config:      config,
		dnsProvider: dnsProvider,
		ipGetter:    ipGetter,
		notifier:    notifier,
	}
}

func (a *App) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := a.updateIP(ctx)
			if err != nil {
				a.logger.Error("update ip error", "error", err.Error())
				os.Exit(1)
			}
			time.Sleep(time.Minute * 3)
		}
	}
}

func (a *App) updateIP(ctx context.Context) error {
	ip, err := a.ipGetter.GetIP()
	if err != nil {
		a.logger.Error("Get ip error", "error", err)
		return err
	}

	name, zone := zoneFromDomain(a.config.Domain)
	a.logger.Info("zone from domain",
		"name", name,
		"zone", zone)

	currentIP, err := a.dnsProvider.GetRecords(ctx, zone)
	if err != nil {
		a.logger.Error("Get records error", "error", err)
		return err
	}

	var found bool
	var record *libdns.Record
	for _, r := range currentIP {
		if r.Name == name {
			found = true
			record = &r
			break
		}
	}
	if found {
		if record.Value == ip {
			a.logger.Info("ip is same, skip update", "ip", ip)
			return nil
		}
		oldIP := record.Value
		record.Value = ip
		_, err = a.dnsProvider.SetRecords(ctx, zone, []libdns.Record{
			*record,
		})
		if err != nil {
			a.logger.Error("Set records error", "error", err)
			return err
		}
		a.notifier.Send(ctx, fmt.Sprintf(IPNotifyFormat, a.config.Domain, oldIP, ip))
	} else {
		_, err = a.dnsProvider.AppendRecords(ctx, zone, []libdns.Record{
			{
				Name:  name,
				Value: ip,
				Type:  "A",
			},
		})
		if err != nil {
			a.logger.Error("Append records error", "error", err)
			return err
		}
	}

	return nil
}

// zoneFromDomain return zone and domain
func zoneFromDomain(domain string) (string, string) {
	arr := strings.SplitN(domain, ".", 2)
	if len(arr) == 1 {
		return "", ""
	}
	return arr[0], arr[1]
}
