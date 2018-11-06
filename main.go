package main

import (
	"fmt"
	"os"
	"time"

	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/notify"
	"github.com/orvice/kit/log"
)

var (
	dnsProvider dns.DNS
	ipGetter    IPGetter

	logger log.Logger
)

func Init() error {
	var err error
	logger = log.NewDefaultLogger()
	GetConfigFromEnv()
	ipGetter = NewIfconfigCo()

	notifier, err := notify.NewTelegramNotifier(TELEGRAM_TOKEN, TELEGRAM_CHATID)
	if err != nil {
		logger.Error(err)
	}

	switch DNS_MODE {
	case DNS_MODE_MU:
		dnsProvider, err = dns.NewMu(API_URI, NODE_ID)
		if err != nil {
			return err
		}
	default:
		dnsProvider, err = dns.NewCloudFlare(CF_API_KEY, CF_API_EMAIL, logger, notifier)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var err error
	err = Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	for {
		updateIP()
		time.Sleep(time.Second * time.Duration(UPDATE_TIME))
	}
}

func updateIP() error {
	ip, err := ipGetter.GetIP()
	if err != nil {
		logger.Errorf("Get ip error %v", err)
		return err
	}
	logger.Infof("Get IP %s", ip)
	err = dnsProvider.UpdateIP(DOMAIN, ip)
	if err != nil {
		logger.Errorf("Update ip error %v", err)
		return err
	}
	return nil
}
