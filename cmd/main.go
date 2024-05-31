package main

import (
	"cloudflare-ddns-go/pkg"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	client, err := pkg.InitClient()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	interval := getIntervalFromEnv("INTERVAL_MINUTES", 5)

	// 立即执行一次 DNS 更新
	updateDNS(client)

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	// 无限循环，保持程序运行
	for {
		select {
		case <-ticker.C:
			updateDNS(client)
		}
	}
}

func getIntervalFromEnv(envVar string, defaultVal int) int {
	valStr := os.Getenv(envVar)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %d minutes. Error: %v", envVar, defaultVal, err)
		return defaultVal
	}

	return val
}

func updateDNS(client *pkg.CloudflareClient) {
	ip, err := pkg.GetExternalIP(client.HTTPClient)
	if err != nil {
		log.Printf("Failed to get external IP: %v", err)
		return
	}

	zoneID, err := client.GetZoneID()
	if err != nil {
		log.Printf("Failed to get zone ID: %v", err)
		return
	}
	log.Printf("Zone ID: %s", zoneID)

	record, err := client.GetRecord(zoneID)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		return
	}

	handleDNSRecord(client, zoneID, record, ip)
	log.Printf("Current DDNS HOST:%s, IP:%s", client.SubdomainName, ip)
}

func handleDNSRecord(client *pkg.CloudflareClient, zoneID string, record *pkg.DNSRecord, ip string) {
	if record == nil {
		log.Printf("Record not found, creating a new one...")
		if recordID, err := client.CreateDNSRecord(zoneID, ip); err != nil {
			log.Printf("Failed to create DNS record: %v", err)
		} else {
			log.Printf("New Record ID: %s", recordID)
		}
	} else if record.Content == ip {
		log.Printf("IP hasn't changed. ip: %s", ip)
	} else {
		if err := client.UpdateDNSRecord(zoneID, record.ID, ip); err != nil {
			log.Printf("Failed to update DNS record: %v", err)
		} else {
			log.Printf("DNS record updated successfully.")
		}
	}
}
