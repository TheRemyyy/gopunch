package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	URLs          []string          `json:"urls"`
	Interval      int               `json:"interval"`
	Timeout       int               `json:"timeout"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers,omitempty"`
	Insecure      bool              `json:"insecure"`
	Follow        bool              `json:"follow_redirects"`
	Concurrency   int               `json:"concurrency"`
	Retries       int               `json:"retries"`
	ExpectedCodes []int             `json:"expected_codes,omitempty"`
	Alerting      *AlertConfig      `json:"alerting,omitempty"`
}

type AlertConfig struct {
	Enabled  bool           `json:"enabled"`
	Cooldown int            `json:"cooldown_seconds"`
	Webhook  *WebhookConfig `json:"webhook,omitempty"`
}

type WebhookConfig struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func parseHeaders(headers []string) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func formatBytes(b int64) string {
	if b < 0 {
		return "-"
	}
	if b < 1024 {
		return fmt.Sprintf("%dB", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(b)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(b)/(1024*1024))
}

// strings package is needed now
