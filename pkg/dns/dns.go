package dns

import (
	"bytes"
	"cloudflare-ddns-go/pkg/client"
	"cloudflare-ddns-go/pkg/env"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ZoneResponse struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DNSRecordResponse struct {
	Result []DNSRecord `json:"result"`
}

func GetZoneID(cfClient *client.CloudflareClient, domain string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfClient.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfClient.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var zoneResp ZoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&zoneResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(zoneResp.Result) == 0 {
		return "", fmt.Errorf("no zones found for domain: %s", domain)
	}

	return zoneResp.Result[0].ID, nil
}

func GetDNSRecord(cfClient *client.CloudflareClient, zoneID, subdomain, domain string) (*DNSRecord, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=A&name=%s.%s", zoneID, subdomain, domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfClient.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfClient.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dnsResp DNSRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&dnsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(dnsResp.Result) == 0 {
		return nil, nil // No record found
	}

	return &dnsResp.Result[0], nil
}

func UpdateDNSRecord(cfClient *client.CloudflareClient, config *env.Config, zoneID, ip string) (string, string, string, error) {
	record, err := GetDNSRecord(cfClient, zoneID, config.Subdomain, config.Domain)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get DNS record: %w", err)
	}

	if record != nil {
		if record.Content == ip {
			return "no_change", ip, ip, nil // No need to update as the IP address is the same
		}
		oldIP := record.Content
		record.Content = ip
		err := putDNSRecord(cfClient, zoneID, record)
		if err != nil {
			return "", "", "", err
		}
		return "updated", oldIP, ip, nil
	}

	// Create new record
	newRecord := &DNSRecord{
		Type:    "A",
		Name:    config.Subdomain + "." + config.Domain,
		Content: ip,
	}
	err = postDNSRecord(cfClient, zoneID, newRecord)
	if err != nil {
		return "", "", "", err
	}
	return "created", "", ip, nil
}

func putDNSRecord(cfClient *client.CloudflareClient, zoneID string, record *DNSRecord) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, record.ID)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal DNS record: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfClient.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfClient.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

func postDNSRecord(cfClient *client.CloudflareClient, zoneID string, record *DNSRecord) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal DNS record: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfClient.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfClient.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}
