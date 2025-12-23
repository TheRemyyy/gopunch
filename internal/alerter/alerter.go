package alerter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Config holds alerting configuration
type Config struct {
	Enabled  bool
	Cooldown time.Duration
	Webhook  *WebhookConfig
}

// WebhookConfig for Discord/Slack/custom webhooks
type WebhookConfig struct {
	URL    string
	Method string
}

// Alert represents an alert event
type Alert struct {
	URL       string
	Status    string
	Error     string
	Timestamp time.Time
}

// Alerter manages sending alerts with cooldown
type Alerter struct {
	config    Config
	lastAlert map[string]time.Time
	mu        sync.Mutex
	client    *http.Client
}

// New creates a new Alerter
func New(config Config) *Alerter {
	return &Alerter{
		config:    config,
		lastAlert: make(map[string]time.Time),
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// SendAlert sends an alert if cooldown has passed
func (a *Alerter) SendAlert(alert Alert) error {
	if !a.config.Enabled {
		return nil
	}

	a.mu.Lock()
	lastTime, exists := a.lastAlert[alert.URL]
	if exists && time.Since(lastTime) < a.config.Cooldown {
		a.mu.Unlock()
		return nil // Still in cooldown
	}
	a.lastAlert[alert.URL] = time.Now()
	a.mu.Unlock()

	if a.config.Webhook != nil {
		return a.sendWebhook(alert)
	}

	return nil
}

func (a *Alerter) sendWebhook(alert Alert) error {
	// Discord webhook format
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "🚨 GoPunch Alert",
				"description": fmt.Sprintf("**URL:** %s\n**Status:** %s", alert.URL, alert.Status),
				"color":       16711680, // Red
				"timestamp":   alert.Timestamp.Format(time.RFC3339),
				"footer": map[string]string{
					"text": "GoPunch Monitoring",
				},
			},
		},
	}

	if alert.Error != "" {
		payload["embeds"].([]map[string]interface{})[0]["description"] =
			fmt.Sprintf("**URL:** %s\n**Error:** %s", alert.URL, alert.Error)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	method := a.config.Webhook.Method
	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequest(method, a.config.Webhook.URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// SendRecoveryAlert sends a recovery notification
func (a *Alerter) SendRecoveryAlert(url string) error {
	if !a.config.Enabled || a.config.Webhook == nil {
		return nil
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "✅ GoPunch Recovery",
				"description": fmt.Sprintf("**URL:** %s\n**Status:** Back online", url),
				"color":       65280, // Green
				"timestamp":   time.Now().Format(time.RFC3339),
				"footer": map[string]string{
					"text": "GoPunch Monitoring",
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", a.config.Webhook.URL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
