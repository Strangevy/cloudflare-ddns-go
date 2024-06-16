package http

import (
	"fmt"
	"io"
	"net/http"
)

const (
	PrimaryURL = "https://api.ipify.org"
	BackupURL  = "https://ifconfig.me"
)

func GetCurrentIP() (string, error) {
	urls := []string{PrimaryURL, BackupURL}

	for _, url := range urls {
		ip, err := fetchIP(url)
		if err != nil {
			fmt.Printf("Failed to fetch IP from %s: %s\n", url, err)
			continue
		}
		return ip, nil
	}

	return "", fmt.Errorf("failed to fetch IP from all provided URLs")
}

func fetchIP(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
