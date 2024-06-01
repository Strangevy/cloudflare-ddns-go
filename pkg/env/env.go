package env

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	APIToken  string
	Domain    string
	Subdomain string
	Interval  int
}

func LoadConfig() *Config {
	apiToken := os.Getenv("CF_API_TOKEN")
	domain := os.Getenv("CF_DOMAIN_NAME")
	subdomain := os.Getenv("CF_SUBDOMAIN_NAME")
	intervalStr := os.Getenv("INTERVAL_MINUTES")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		interval = 5 // Default to 5 minutes if interval is not set or invalid
	}

	if apiToken == "" || domain == "" || subdomain == "" {
		log.Println("Missing required environment variables")
		return nil
	}

	return &Config{
		APIToken:  apiToken,
		Domain:    domain,
		Subdomain: subdomain,
		Interval:  interval,
	}
}
