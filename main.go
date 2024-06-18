package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	SelectedProvider := Cloudflare // default provider
	flag.Var(&SelectedProvider, "provider", "the provider for manage DNS record")

	var creates, deletes stringSlice
	flag.Var(&creates, "create", "specify a IP to create, can define multiple times")
	flag.Var(&deletes, "delete", "specify a IP to delete, can define multiple times")

	var domain string
	flag.StringVar(&domain, "domain", "", "specify the domain name")

	flag.Parse()
	slog.Info(fmt.Sprintf("using %s as DNS provider", SelectedProvider))

	if len(creates) == 0 && len(deletes) == 0 {
		fmt.Println("-create or -delete option should be defined")
		os.Exit(1)
	}

	if domain == "" {
		fmt.Println("please define domain name")
		os.Exit(1)
	}

	ctx := context.Background()

	api, err := NewDNSProvider(SelectedProvider)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to initialize cloudflare provider: %s", err.Error()))
	}

	for _, c := range creates {
		if err := api.Create(ctx, domain, c); err != nil {
			slog.Error(err.Error())
		}
	}

	for _, d := range deletes {
		if err := api.Delete(ctx, domain, d); err != nil {
			slog.Error(err.Error())
		}
	}

}
