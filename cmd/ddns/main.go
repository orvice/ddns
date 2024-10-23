package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/weeon/log"
	"github.com/weeon/utils/task"

	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/config"
	"github.com/orvice/ddns/internal/ip"
	"github.com/orvice/ddns/notify"
)

var (
	dnsProvider dns.LibDNS
	ipGetter    ip.Getter
	DNSMode     string
)

func Init() error {
	var err error
	config.GetConfigFromEnv()
	ipGetter = ip.NewIfconfigCo()

	notify.Init()

	notifier, err := notify.NewTelegramNotifier(config.TelegramToken, config.TelegramChatID)
	if err != nil {
		log.Errorf("notify init error %v", err)
	} else {
		notify.AddNotifier(notifier)
	}

	switch config.DNSMode {
	default:
		dnsProvider = dns.NewCloudFlare()
	}

	return nil
}

func main() {
	var err error
	log.FastInitFileLogger()
	err = Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	ctx := context.Background()

	task.NewTaskAndRun("updateUpdate", time.Minute*3, func() error {
		err := updateIP(ctx)
		if err != nil {
			slog.Error("update ip error", "error", err.Error())
			os.Exit(1)
		}
		return err
	}, task.SetTaskLogger(log.GetDefault()))
	select {}
}

func updateIP(ctx context.Context) error {
	logger := slog.Default()
	ip, err := ipGetter.GetIP()
	if err != nil {
		logger.Error("Get ip error", "error", err)
		return err
	}

	name, zone := zoneFromDomain(config.DOMAIN)
	logger.Info("zone from domain",
		"name", name,
		"zone", zone)

	currentIP, err := dnsProvider.GetRecords(ctx, zone)
	if err != nil {
		logger.Error("Get records error", "error", err)
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
			logger.Info("ip is same, skip update", "ip", ip)
			return nil
		}
		record.Value = ip
		_, err = dnsProvider.SetRecords(ctx, zone, []libdns.Record{
			*record,
		})
		if err != nil {
			logger.Error("Set records error", "error", err)
			return err
		}
	} else {
		_, err = dnsProvider.AppendRecords(ctx, zone, []libdns.Record{
			{
				Name:  name,
				Value: ip,
			},
		})
		if err != nil {
			logger.Error("Append records error", "error", err)
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
