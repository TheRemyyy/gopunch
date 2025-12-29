# Alerting System

GoPunch features a built-in alerting engine designed to notify your team via Discord or Slack webhooks when service status changes.

## How it Works

Alerts are triggered during `watch` mode. The engine tracks the state of each URL and fires notifications based on transitions:

1.  **Healthy -> Unhealthy**: A "Failure" alert is sent immediately.
2.  **Unhealthy -> Healthy**: A "Recovery" alert is sent once the service is back online.

## Configuration

Alerts are configured in the `alerting` section of `gopunch.json`:

```json
"alerting": {
  "enabled": true,
  "cooldown_seconds": 300,
  "webhook": {
    "url": "https://discord.com/api/webhooks/...",
    "method": "POST"
  }
}
```

### Parameters

- **enabled**: Set to `true` to activate the alerting engine.
- **cooldown_seconds**: Prevents notification spam. If a service stays down, GoPunch will wait this many seconds before sending another "Failure" alert. Default is 300 seconds (5 minutes).
- **webhook.url**: The full URL of your Discord or Slack webhook.
- **webhook.method**: HTTP method for the webhook call (usually `POST`).

## Discord Integration

GoPunch sends rich embeds to Discord for maximum readability:

- 🔴 **Failure Alerts**: Include the URL, status code (if available), error message, and timestamp.
- 🟢 **Recovery Alerts**: Clear notification that the service is "Back online".

## Slack Integration

Since GoPunch uses a standard JSON payload, it can be easily adapted for Slack by using their "Incoming Webhooks" feature.

## Cooldown Logic

The cooldown is per-URL. If `https://a.com` and `https://b.com` both go down, you will receive two separate alerts. Subsequent failures for the same URL will be suppressed until the cooldown timer expires, at which point one fresh alert will be sent if the service is still down.
