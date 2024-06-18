package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareDNS struct {
	api    *cloudflare.API
	zoneId string
}

func NewCloudflareDNS(api *cloudflare.API, zoneId string) *CloudflareDNS {
	return &CloudflareDNS{
		api:    api,
		zoneId: zoneId,
	}
}

func (c *CloudflareDNS) GetZoneDetails(ctx context.Context) (*cloudflare.Zone, error) {
	zoneData, err := c.api.ZoneDetails(ctx, c.zoneId)
	if err != nil {
		return nil, err
	}

	return &zoneData, nil
}

func (c *CloudflareDNS) GetListRecords(ctx context.Context, name string, IP string) ([]cloudflare.DNSRecord, error) {

	dnsRecords, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(c.zoneId), cloudflare.ListDNSRecordsParams{
		Name:    name,
		Type:    "A",
		Content: IP,
	})

	if err != nil {
		slog.Error(fmt.Sprintf("can't get list of dns record for %s with IP %s: %s", name, IP, err.Error()))
		return nil, err
	}

	return dnsRecords, nil
}

func (c *CloudflareDNS) Create(ctx context.Context, name string, IP string) error {

	record := cloudflare.CreateDNSRecordParams{
		Name:    name,
		Type:    "A",
		Content: IP,
		Proxied: cloudflare.BoolPtr(false),
	}

	_, err := c.api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(c.zoneId), record)
	if err != nil {
		slog.Error(fmt.Sprintf("can't create record for %s with IP %s: %s", name, IP, err.Error()))
		return err
	}

	slog.Info("record created")
	return nil
}

func (c *CloudflareDNS) Delete(ctx context.Context, name string, IP string) error {

	zoneDetail, err := c.GetZoneDetails(ctx)
	if err != nil {
		return err
	}

	dnsRecords, err := c.GetListRecords(ctx, name+"."+zoneDetail.Name, IP)
	if err != nil {
		return err
	}

	if len(dnsRecords) <= 0 {
		return fmt.Errorf("DNS record %s with IP %s not found in record", name, IP)
	}

	firstRecord := dnsRecords[0]

	if err := c.api.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(c.zoneId), firstRecord.ID); err != nil {
		slog.Error(fmt.Sprintf("can't delete DNS record for %s with IP %s and record.ID %s: %s", name, IP, firstRecord.ID, err.Error()))
		return err
	}

	slog.Info("record deleted")
	return nil
}
