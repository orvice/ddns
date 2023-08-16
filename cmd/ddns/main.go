package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/catpie/musdk-go"
	"github.com/weeon/utils/task"

	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/config"
	"github.com/orvice/ddns/internal/ip"
	"github.com/orvice/ddns/notify"
	"github.com/weeon/log"
)

var (
	dnsProvider dns.DNS
	ipGetter    ip.IPGetter
	muCli       *musdk.Client
)

func Init() error {

	var err error
	config.GetConfigFromEnv()
	ipGetter = ip.NewIfconfigCo()
	muCli = musdk.ClientFromEnv(slog.Default())

	notify.Init()

	notifier, err := notify.NewTelegramNotifier(config.TELEGRAM_TOKEN, config.TELEGRAM_CHATID)
	if err != nil {
		log.Errorf("notify init error %v", err)
	} else {
		notify.AddNotifier(notifier)
	}

	switch config.DNS_MODE {
	case config.DNS_MODE_MU:
		dnsProvider, err = dns.NewMu(config.API_URI, config.NODE_ID)
		if err != nil {
			return err
		}
	default:
		dnsProvider, err = dns.NewCloudFlare(config.CF_API_KEY, config.CF_API_EMAIL, log.GetDefault())
		if err != nil {
			return err
		}
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
		return updateIP(ctx)
	}, task.SetTaskLogger(log.GetDefault()))
	select {}
}

func updateIP(ctx context.Context) error {
	ip, err := ipGetter.GetIP()
	if err != nil {
		log.Errorf("Get ip error %v", err)
		return err
	}

	currentIP, err := dnsProvider.GetIP(ctx, config.DOMAIN)

	log.Infow("get ip",
		"from_getter", ip,
		"from_provider", currentIP)
	err = dnsProvider.UpdateIP(ctx, config.DOMAIN, ip)
	if err != nil {
		log.Errorf("Update ip error %v", err)
		return err
	}
	return nil
}
