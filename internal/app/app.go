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
	// Attach domain to every log line so operators always know which target
	// produced the entry.
	logger = logger.With("domain", config.Domain)
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
	name, zone := zoneFromDomain(a.config.Domain)
	a.logger.Info("ddns started",
		"interval", a.interval,
		"dns_provider", a.config.DNSProvider,
		"name", name,
		"zone", zone,
	)

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	if err := a.UpdateIP(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			a.logger.Info("ddns stopping", "reason", ctx.Err())
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
		a.logger.Error("get ip failed", "error", err)
		return err
	}

	name, zone := zoneFromDomain(a.config.Domain)
	a.logger.Debug("resolved zone", "name", name, "zone", zone)

	records, err := a.dnsProvider.GetRecords(ctx, zone)
	if err != nil {
		a.logger.Error("get records failed", "zone", zone, "error", err)
		return err
	}

	record, found := findRecord(records, name)
	if !found {
		return a.appendRecord(ctx, zone, name, currentIP)
	}
	if record.Value == currentIP {
		a.logger.Info("ip unchanged, skip update", "ip", currentIP)
		return nil
	}
	return a.updateRecord(ctx, zone, record, currentIP)
}

func (a *App) appendRecord(ctx context.Context, zone, name, ip string) error {
	_, err := a.dnsProvider.AppendRecords(ctx, zone, []libdns.Record{
		{Name: name, Value: ip, Type: "A"},
	})
	if err != nil {
		a.logger.Error("append record failed", "zone", zone, "name", name, "ip", ip, "error", err)
		return err
	}
	a.logger.Info("record appended", "zone", zone, "name", name, "ip", ip)
	return nil
}

func (a *App) updateRecord(ctx context.Context, zone string, record libdns.Record, newIP string) error {
	oldIP := record.Value
	record.Value = newIP
	if _, err := a.dnsProvider.SetRecords(ctx, zone, []libdns.Record{record}); err != nil {
		a.logger.Error("set record failed", "zone", zone, "name", record.Name, "old_ip", oldIP, "new_ip", newIP, "error", err)
		return err
	}
	a.logger.Info("record updated", "zone", zone, "name", record.Name, "old_ip", oldIP, "new_ip", newIP)
	if a.notifier != nil {
		if err := a.notifier.Send(ctx, fmt.Sprintf(IPNotifyFormat, a.config.Domain, oldIP, newIP)); err != nil {
			a.logger.Warn("notify failed", "error", err)
		}
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
