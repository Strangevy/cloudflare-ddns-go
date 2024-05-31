package pkg

import "net/http"

type CloudflareClient struct {
	APIToken      string
	DomainName    string
	SubdomainName string
	HTTPClient    *http.Client
}

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}
