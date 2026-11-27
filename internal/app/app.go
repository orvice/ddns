package app

import (
	"context"
	"fmt"
	"log/slog"
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

	DefaultInterval = 3 * time.Minute
)

type App struct {
	logger      *slog.Logger
	config      *config.Config
	dnsProvider dns.LibDNS
	ipGetter    ip.Getter
	notifier    notify.Notifier
	interval    time.Duration
}

func New(logger *slog.Logger, config *config.Config, dnsProvider dns.LibDNS, ipGetter ip.Getter, notifier notify.Notifier) *App {
	return &App{
		logger:      logger,
		config:      config,
		dnsProvider: dnsProvider,
		ipGetter:    ipGetter,
		notifier:    notifier,
		interval:    DefaultInterval,
	}
}

// WithInterval overrides the polling interval. Useful for tests.
func (a *App) WithInterval(d time.Duration) *App {
	a.interval = d
	return a
}

// Run polls and syncs until ctx is cancelled or an unrecoverable error occurs.
// It returns the cause; nil if ctx was cancelled cleanly.
func (a *App) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	if err := a.UpdateIP(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := a.UpdateIP(ctx); err != nil {
				return err
			}
		}
	}
}

// UpdateIP performs one sync cycle: fetch current IP, look up the existing
// record, and create or update it as needed.
func (a *App) UpdateIP(ctx context.Context) error {
	currentIP, err := a.ipGetter.GetIP()
	if err != nil {
		a.logger.Error("Get ip error", "error", err)
		return err
	}

	name, zone := zoneFromDomain(a.config.Domain)
	a.logger.Info("zone from domain", "name", name, "zone", zone)

	records, err := a.dnsProvider.GetRecords(ctx, zone)
	if err != nil {
		a.logger.Error("Get records error", "error", err)
		return err
	}

	record, found := findRecord(records, name)
	if !found {
		return a.appendRecord(ctx, zone, name, currentIP)
	}
	if record.Value == currentIP {
		a.logger.Info("ip is same, skip update", "ip", currentIP)
		return nil
	}
	return a.updateRecord(ctx, zone, record, currentIP)
}

func (a *App) appendRecord(ctx context.Context, zone, name, ip string) error {
	_, err := a.dnsProvider.AppendRecords(ctx, zone, []libdns.Record{
		{Name: name, Value: ip, Type: "A"},
	})
	if err != nil {
		a.logger.Error("Append records error", "error", err)
	}
	return err
}

func (a *App) updateRecord(ctx context.Context, zone string, record libdns.Record, newIP string) error {
	oldIP := record.Value
	record.Value = newIP
	if _, err := a.dnsProvider.SetRecords(ctx, zone, []libdns.Record{record}); err != nil {
		a.logger.Error("Set records error", "error", err)
		return err
	}
	if a.notifier != nil {
		_ = a.notifier.Send(ctx, fmt.Sprintf(IPNotifyFormat, a.config.Domain, oldIP, newIP))
	}
	return nil
}

func findRecord(records []libdns.Record, name string) (libdns.Record, bool) {
	for _, r := range records {
		if r.Name == name {
			return r, true
		}
	}
	return libdns.Record{}, false
}

// zoneFromDomain return name and zone.
func zoneFromDomain(domain string) (string, string) {
	arr := strings.SplitN(domain, ".", 2)
	if len(arr) == 1 {
		return "", ""
	}
	return arr[0], arr[1]
}
