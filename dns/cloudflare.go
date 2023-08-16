package dns

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/orvice/ddns/config"
	"github.com/orvice/ddns/notify"
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

	token := os.Getenv("CF_TOKEN")
	if token != "" {
		slog.Info("cf use token")
		client, err = cloudflare.NewWithAPIToken(token)
		if err != nil {
			return nil, err
		}
	}

	return &CloudFlare{
		client: client,
		logger: logger,
	}, nil
}

func (c *CloudFlare) getZone(ctx context.Context, domain string) (*cloudflare.Zone, error) {
	zones, err := c.client.ListZones(ctx)
	if err != nil {
		slog.Error("list zones error", "err", err)
		return nil, err
	}

	for _, z := range zones {
		if strings.Contains(domain, z.Name) {
			slog.Info("get zone success", "zone", z.Name, "id", z.ID)
			return &z, nil
		}
	}
	return nil, fmt.Errorf("not found zone")
}

func (c *CloudFlare) GetDomainZoneID(domain string) (string, error) {
	var zoneID = os.Getenv("CF_ZONE_ID")
	if zoneID != "" {
		return zoneID, nil
	}
	zone, err := c.getZone(context.Background(), domain)
	if err != nil {
		slog.Error("get zone error", "err", err)
	}
	return zone.ID, nil
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
				slog.Error("update dns record error",
					"err", err)
				return err
			}
			notify.Notify(ctx, fmt.Sprintf(config.IpNotifyFormat, domain, oldIP, ip))
		}
	}
	return nil
}
