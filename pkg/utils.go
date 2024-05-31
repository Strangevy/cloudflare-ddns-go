package pkg

import (
	"fmt"
	"io"
	"net/http"
)

// 获取当前的外部IP地址
func GetExternalIP(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
