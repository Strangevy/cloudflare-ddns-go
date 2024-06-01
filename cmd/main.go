package main

import (
	"cloudflare-ddns-go/pkg/client"
	"cloudflare-ddns-go/pkg/dns"
	"cloudflare-ddns-go/pkg/env"
	"cloudflare-ddns-go/pkg/http"
	"log"
	"time"
)

var cachedZoneID string

func main() {
	config := env.LoadConfig()
	if config == nil {
		log.Fatal("Failed to load configuration")
	}

	cfClient := client.NewCloudflareClient(config)

	for {
		currentIP, err := http.GetCurrentIP()
		if err != nil {
			log.Fatalf("Failed to get current IP: %v", err)
		}

		if cachedZoneID == "" {
			cachedZoneID, err = dns.GetZoneID(cfClient, config.Domain)
			if err != nil {
				log.Fatalf("Failed to get Zone ID: %v", err)
			}
		}

		status, oldIP, newIP, err := dns.UpdateDNSRecord(cfClient, config, cachedZoneID, currentIP)
		if err != nil {
			log.Fatalf("Failed to update DNS record: %v", err)
		}

		if status == "no_change" {
			log.Printf("No change needed. Current IP: %s, Domain: %s.%s\n", currentIP, config.Subdomain, config.Domain)
		} else if status == "updated" {
			log.Printf("Updated DNS record. Domain: %s.%s, Old IP: %s, New IP: %s\n", config.Subdomain, config.Domain, oldIP, newIP)
		} else if status == "created" {
			log.Printf("Created new DNS record. Domain: %s.%s, IP: %s\n", config.Subdomain, config.Domain, newIP)
		}

		time.Sleep(time.Duration(config.Interval) * time.Minute)
	}
}
