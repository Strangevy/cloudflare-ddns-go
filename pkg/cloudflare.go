package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL        = "https://api.cloudflare.com/client/v4"
	httpMethodGet  = "GET"
	httpMethodPost = "POST"
	httpMethodPut  = "PUT"
)

type ListResponse struct {
	Success bool        `json:"success"`
	Result  []DNSRecord `json:"result"`
}

type ModifyResponse struct {
	Success bool        `json:"success"`
	Result  DNSRecord `json:"result"`
}

// 创建和发送HTTP请求，并解码响应
func (c *CloudflareClient) doRequest(method, url string, body []byte, target interface{}) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// 获取 Zone ID
func (c *CloudflareClient) GetZoneID() (string, error) {
	url := fmt.Sprintf("%s/zones?name=%s", baseURL, c.DomainName)
	var result ListResponse
	if err := c.doRequest(httpMethodGet, url, nil, &result); err != nil {
		return "", err
	}

	if result.Success && len(result.Result) > 0 {
		return result.Result[0].ID, nil
	}

	return "", fmt.Errorf("unable to get zone ID")
}

// 获取 Record 对象
func (c *CloudflareClient) GetRecord(zoneID string) (*DNSRecord, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records?name=%s", baseURL, zoneID, c.SubdomainName)
	var result ListResponse
	if err := c.doRequest(httpMethodGet, url, nil, &result); err != nil {
		return nil, err
	}

	if result.Success && len(result.Result) > 0 {
		return &result.Result[0], nil
	}else if result.Success{
		return nil, nil
	}

	return nil, fmt.Errorf("unable to get record")
}

// 创建新的 DNS 记录
func (c *CloudflareClient) CreateDNSRecord(zoneID, ipAddress string) (string, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records", baseURL, zoneID)
	record := DNSRecord{
		Type:    "A",
		Name:    c.SubdomainName,
		Content: ipAddress,
		TTL:     120,
		Proxied: false,
	}
	jsonData, err := json.Marshal(record)
	if err != nil {
		return "", err
	}

	var result ModifyResponse
	if err := c.doRequest(httpMethodPost, url, jsonData, &result); err != nil {
		return "", err
	}

	if result.Success {
		return result.Result.ID, nil
	}

	return "", fmt.Errorf("failed to create DNS record")
}

// 更新 DNS 记录
func (c *CloudflareClient) UpdateDNSRecord(zoneID, recordID, ipAddress string) error {
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", baseURL, zoneID, recordID)
	record := DNSRecord{
		Type:    "A",
		Name:    c.SubdomainName,
		Content: ipAddress,
		TTL:     120,
		Proxied: false,
	}
	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	var result ModifyResponse
	if err := c.doRequest(httpMethodPut, url, jsonData, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to update DNS record")
	}

	return nil
}
