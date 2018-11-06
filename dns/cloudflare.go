package dns

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/orvice/ddns/notify"
	"github.com/orvice/ddns/utils"
	"github.com/orvice/kit/log"
)

type CloudFlare struct {
	client *cloudflare.API
	logger log.Logger
	notify notify.Notifier
}

func NewCloudFlare(key, email string, logger log.Logger, notify notify.Notifier) (*CloudFlare, error) {
	client, err := cloudflare.New(key, email)
	if err != nil {
		return nil, err
	}
	return &CloudFlare{
		client: client,
		logger: logger,
		notify: notify,
	}, nil
}

func (c *CloudFlare) GetDomainZoneID(domain string) (string, error) {
	zone := utils.GetDomainSuffix(domain)
	id, err := c.client.ZoneIDByName(zone)
	if err != nil {
		return "", nil
	}
	return id, nil
}

// @todo
func (c *CloudFlare) GetIP(domain string) (string, error) {
	return "", nil
}

func (c *CloudFlare) UpdateIP(domain, ip string) error {
	zid, err := c.GetDomainZoneID(domain)
	if err != nil {
		return err
	}
	rs, err := c.client.DNSRecords(zid, cloudflare.DNSRecord{
		Name: domain,
	})
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Type == "A" {
			if r.Content == ip {
				c.logger.Infof("ip not change...")
				continue
			}
			oldIP := r.Content
			r.Content = ip
			err = c.client.UpdateDNSRecord(zid, r.ID, r)
			if err != nil {
				return err
			}
			c.logger.Infof("ip update success...")
			if c.notify != nil {
				err = c.notify.Send(fmt.Sprintf("[%s] ip changed, old IP: %s new IP: %s",
					domain, oldIP, ip))
				if err != nil {
					c.logger.Errorf("notify error %s", err.Error())
				}
			} else {
				c.logger.Infof("skip notify")
			}
		}
	}
	return nil
}
