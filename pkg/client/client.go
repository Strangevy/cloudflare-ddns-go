package client

import (
	"cloudflare-ddns-go/pkg/env"
	"net/http"
)

type CloudflareClient struct {
	APIToken string
	Client   *http.Client
}

func NewCloudflareClient(config *env.Config) *CloudflareClient {
	return &CloudflareClient{
		APIToken: config.APIToken,
		Client:   &http.Client{},
	}
}
