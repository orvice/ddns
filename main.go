package main

import (
	"fmt"
	"github.com/orvice/ddns/dns"
	"os"
	"time"
)

var (
	dnsProvider dns.DNS
	ipGetter    IPGetter
)

func Init() error {
	var err error
	GetConfigFromEnv()
	ipGetter = NewIfconfigCo()
	dnsProvider, err = dns.NewCloudFlare(CF_API_KEY, CF_API_EMAIL)
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
		return err
	}
	err = dnsProvider.UpdateIP(DOMAIN, ip)
	if err != nil {
		return err
	}
	return nil
}
