package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type Provider string

const (
	AWS        Provider = "aws"
	Cloudflare Provider = "cloudflare"
)

// implement flag.Value
func (p *Provider) Set(value string) error {
	switch value {
	case string(AWS), string(Cloudflare):
		*p = Provider(value)
		return nil
	default:
		return errors.New("invalid provider, must be 'aws' or 'cloudflare'")
	}
}

func (p *Provider) String() string {
	return string(*p)
}

type DNSProvider interface {
	Create(ctx context.Context, name string, IP string) error
	Delete(ctx context.Context, name string, IP string) error
}

func NewDNSProvider(provider Provider) (DNSProvider, error) {
	switch provider {
	case Cloudflare:
		cfZoneId := os.Getenv("CF_ZONEID")
		cfToken := os.Getenv("CF_TOKEN")

		if cfZoneId == "" || cfToken == "" {
			fmt.Println("Please define CF_ZONEID & CF_TOKEN")
			os.Exit(1)
		}

		cfApi, err := cloudflare.NewWithAPIToken(cfToken)
		if err != nil {
			return nil, err
		}

		return NewCloudflareDNS(cfApi, cfZoneId), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
