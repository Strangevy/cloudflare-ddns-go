package main

import (
	"cloudflare-ddns-go/pkg"
	"log"
)

func main() {
	client, err := pkg.InitClient()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	ip, err := pkg.GetExternalIP(client.HTTPClient)
	if err != nil {
		log.Fatalf("Failed to get external IP: %v", err)
	}

	zoneID, err := client.GetZoneID()
	if err != nil {
		log.Fatalf("Failed to get zone ID: %v", err)
	}
	log.Printf("Zone ID: %s", zoneID)

	record, err := client.GetRecord(zoneID)
	if err != nil {
		log.Fatalf("Failed to get record: %v", err)
	}
	if record == nil {
		log.Printf("Record not found, creating a new one...")
		recordID, err := client.CreateDNSRecord(zoneID, ip)
		if err != nil {
			log.Fatalf("Failed to create DNS record: %v", err)
		}
		log.Printf("New Record ID: %s", recordID)
	} else if record.Content == ip {
		log.Printf("IP hasn't changed. ip: %s", ip)
	} else {
		if err := client.UpdateDNSRecord(zoneID, record.ID, ip); err != nil {
			log.Fatalf("Failed to update DNS record: %v", err)
		}
		log.Printf("DNS record updated successfully.")
	}
	log.Printf("Current DDNS HOST:%s, IP:%s", client.SubdomainName, ip)
}
