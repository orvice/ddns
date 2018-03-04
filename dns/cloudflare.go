package dns

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/orvice/ddns/utils"
)

type CloudFlare struct {
	client *cloudflare.API
}

func NewCloudFlare(key, email string) (*CloudFlare, error) {
	client, err := cloudflare.New(key, email)
	if err != nil {
		return nil, err
	}
	return &CloudFlare{
		client: client,
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
				continue
			}
			r.Content = ip
			err = c.client.UpdateDNSRecord(zid, r.ID, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
