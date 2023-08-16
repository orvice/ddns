package dns

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/orvice/ddns/config"
	"github.com/orvice/ddns/notify"
	"github.com/orvice/ddns/utils"
	"github.com/weeon/contract"
)

type CloudFlare struct {
	client *cloudflare.API
	logger contract.Logger
}

func NewCloudFlare(key, email string, logger contract.Logger) (*CloudFlare, error) {
	client, err := cloudflare.New(key, email)
	if err != nil {
		return nil, err
	}
	return &CloudFlare{
		client: client,
		logger: logger,
	}, nil
}

func (c *CloudFlare) GetDomainZoneID(domain string) (string, error) {
	var domainID = os.Getenv("CF_DOMAIN_ID")
	if domainID != "" {
		return domainID, nil
	}
	zone := utils.GetDomainSuffix(domain)
	id, err := c.client.ZoneIDByName(zone)
	if err != nil {
		return "", nil
	}
	return id, nil
}

func (c *CloudFlare) GetIP(ctx context.Context, domain string) (string, error) {
	zid, err := c.GetDomainZoneID(domain)
	if err != nil {
		return "", err
	}
	rs, ret, err := c.client.ListDNSRecords(ctx, &cloudflare.ResourceContainer{
		// Name: domain,
		Identifier: zid,
	}, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return "", err
	}
	slog.Info("get dns records ", "count", ret.Count)

	for _, v := range rs {
		return v.Content, nil
	}

	return "", nil
}

func (c *CloudFlare) UpdateIP(ctx context.Context, domain, ip string) error {
	zid, err := c.GetDomainZoneID(domain)
	if err != nil {
		return err
	}
	rs, ret, err := c.client.ListDNSRecords(ctx, &cloudflare.ResourceContainer{
		// Name: domain,
		Identifier: zid,
	}, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return err
	}
	slog.Info("get dns records ", "count", ret.Count)
	for _, r := range rs {
		if r.Type == "A" {
			if r.Content == ip {
				c.logger.Infof("ip not change...")
				continue
			}
			oldIP := r.Content
			r.Content = ip
			r, err = c.client.UpdateDNSRecord(ctx, &cloudflare.ResourceContainer{
				Identifier: r.ID,
			}, cloudflare.UpdateDNSRecordParams{
				Type:    r.Type,
				Name:    r.Name,
				Content: ip,
			})

			slog.Info("update dns record",
				"record_id", r.ID,
				"domain", domain, "old_ip", oldIP, "new_ip", ip)

			if err != nil {
				return err
			}
			notify.Notify(ctx, fmt.Sprintf(config.IpNotifyFormat, domain, oldIP, ip))
		}
	}
	return nil
}
