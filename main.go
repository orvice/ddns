package main

import (
	"fmt"
	"github.com/orvice/ddns/dns"
	"github.com/orvice/kit/log"
	"os"
	"time"
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
	dnsProvider, err = dns.NewCloudFlare(CF_API_KEY, CF_API_EMAIL, logger)
	if err != nil {
		return err
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
