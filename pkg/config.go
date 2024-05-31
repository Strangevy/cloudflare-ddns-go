package pkg

import (
	"fmt"
	"net/http"
	"os"
)

func InitClient() (*CloudflareClient, error) {
	apiToken := os.Getenv("CF_API_TOKEN")
	domainName := os.Getenv("CF_DOMAIN_NAME")
	subdomainName := os.Getenv("CF_SUBDOMAIN_NAME")

	if apiToken == "" || domainName == "" || subdomainName == "" {
		return nil, fmt.Errorf("one or more required environment variables are not set")
	}

	client := &CloudflareClient{
		APIToken:      apiToken,
		DomainName:    domainName,
		SubdomainName: subdomainName,
		HTTPClient:    &http.Client{},
	}

	return client, nil
}
