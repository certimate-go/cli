# Notification Provider Configuration Reference

This document provides detailed configuration for all notification providers. Each provider has specific access config fields (stored in the access record) and optional extended config fields (in the workflow node's `providerConfig`).

> **Important**: The field name for provider type is `provider` (not `providerType`), and the provider-specific config field is `providerConfig` (not `extendedConfig`).

## How Notification Configuration Works

In the workflow `graphDraft`, a `bizNotify` node has this structure:

```json
{
  "id": "notify-node",
  "type": "bizNotify",
  "data": {
    "name": "Send Notification",
    "config": {
      "provider": "email",
      "providerAccessId": "ACCESS_ID",
      "providerConfig": {
        "receiverAddress": "admin@example.com",
        "format": "html"
      },
      "subject": "[Certimate] Workflow Alert",
      "message": "Workflow completed successfully.",
      "skipOnAllPrevSkipped": false
    }
  }
}
```

**Key Fields:**
- `provider`: The notification provider type (see below)
- `providerAccessId`: ID of the access credential (created via `certimate access create`)
- `providerConfig`: Provider-specific extended configuration (optional for most providers)
- `subject`: Notification subject (supports `{{$variable}}` templates)
- `message`: Notification body (supports `{{$variable}}` templates)
- `skipOnAllPrevSkipped`: Skip if all previous nodes were skipped

---

## Email (`email`)

Sends notifications via SMTP email.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `smtpHost` | string | Yes | SMTP server address (e.g., `smtp.gmail.com`) |
| `smtpPort` | int32 | Yes | SMTP server port (0 = auto based on TLS setting) |
| `smtpTls` | bool | Yes | Whether to enable TLS |
| `username` | string | Yes | SMTP username |
| `password` | string | Yes | SMTP password |
| `senderAddress` | string | Yes | Sender email address |
| `senderName` | string | No | Sender display name |
| `receiverAddress` | string | No* | Default receiver email address |
| `allowInsecureConnections` | bool | No | Allow insecure connections |

### ProviderConfig (in workflow node)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `receiverAddress` | string | Yes* | Receiver email address (fallback if not in access) |
| `format` | string | No | Message format: `"plain"` or `"html"` (default: `"plain"`) |

*`receiverAddress` must be provided in either access config or providerConfig

**Example:**
```json
{
  "provider": "email",
  "providerAccessId": "ACCESS_EMAIL",
  "providerConfig": {
    "receiverAddress": "admin@example.com",
    "format": "html"
  },
  "subject": "[Certimate] Certificate Renewal",
  "message": "<h2>Certificate Renewed</h2><p>Domain: {{$certificate.subjectAltNames}}</p>"
}
```

---

## DingTalk Bot (`dingtalkbot`)

Sends notifications via DingTalk robot webhook.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `webhookUrl` | string | Yes | DingTalk robot webhook URL |
| `secret` | string | No | DingTalk robot secret for signing (recommended) |
| `customPayload` | string | No | Custom JSON message payload |

### ProviderConfig (in workflow node)

None required. Subject and message are sent automatically.

**Example:**
```json
{
  "provider": "dingtalkbot",
  "providerAccessId": "ACCESS_DINGTALK",
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Lark Bot (`larkbot`) / Feishu Bot

Sends notifications via Lark/Feishu robot webhook.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `webhookUrl` | string | Yes | Lark/Feishu robot webhook URL |
| `secret` | string | No | Lark robot secret for signing (recommended) |
| `customPayload` | string | No | Custom JSON message payload |

### ProviderConfig (in workflow node)

None required.

**Example:**
```json
{
  "provider": "larkbot",
  "providerAccessId": "ACCESS_LARK",
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## WeCom Bot (`wecombot`) / Enterprise WeChat Bot

Sends notifications via WeCom robot webhook.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `webhookUrl` | string | Yes | WeCom robot webhook URL |
| `customPayload` | string | No | Custom JSON message payload |

### ProviderConfig (in workflow node)

None required.

**Example:**
```json
{
  "provider": "wecombot",
  "providerAccessId": "ACCESS_WECOM",
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Slack Bot (`slackbot`)

Sends notifications via Slack Bot API.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `botToken` | string | Yes | Slack Bot API Token (starts with `xoxb-`) |
| `channelId` | string | No* | Slack Channel ID |

### ProviderConfig (in workflow node)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `channelId` | string | Yes* | Slack Channel ID (fallback if not in access) |

*`channelId` must be provided in either access config or providerConfig

**Example:**
```json
{
  "provider": "slackbot",
  "providerAccessId": "ACCESS_SLACK",
  "providerConfig": {
    "channelId": "C12345678"
  },
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Discord Bot (`discordbot`)

Sends notifications via Discord Bot API.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `botToken` | string | Yes | Discord Bot API Token |
| `channelId` | string | No* | Discord Channel ID |

### ProviderConfig (in workflow node)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `channelId` | string | Yes* | Discord Channel ID (fallback if not in access) |

*`channelId` must be provided in either access config or providerConfig

**Example:**
```json
{
  "provider": "discordbot",
  "providerAccessId": "ACCESS_DISCORD",
  "providerConfig": {
    "channelId": "123456789012345678"
  },
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Telegram Bot (`telegrambot`)

Sends notifications via Telegram Bot API.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `botToken` | string | Yes | Telegram Bot API Token (from @BotFather) |
| `chatId` | string | No* | Telegram Chat ID |

### ProviderConfig (in workflow node)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `chatId` | string | Yes* | Telegram Chat ID (fallback if not in access) |

*`chatId` must be provided in either access config or providerConfig

**Example:**
```json
{
  "provider": "telegrambot",
  "providerAccessId": "ACCESS_TELEGRAM",
  "providerConfig": {
    "chatId": "-1001234567890"
  },
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Mattermost (`mattermost`)

Sends notifications via Mattermost API.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `serverUrl` | string | Yes | Mattermost server URL |
| `username` | string | Yes | Mattermost username |
| `password` | string | Yes | Mattermost password |
| `channelId` | string | No* | Mattermost channel ID |

### ProviderConfig (in workflow node)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `channelId` | string | Yes* | Mattermost channel ID (fallback if not in access) |

*`channelId` must be provided in either access config or providerConfig

**Example:**
```json
{
  "provider": "mattermost",
  "providerAccessId": "ACCESS_MATTERMOST",
  "providerConfig": {
    "channelId": "channel-id-here"
  },
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Webhook (`webhook`)

Sends notifications via custom HTTP webhook.

### Access Config (stored in access record)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | Yes | Webhook URL |
| `method` | string | No | HTTP method (default: `"POST"`) |
| `headers` | string | No | Request headers (string format) |
| `data` | string | No | Default request body data (JSON string) |
| `allowInsecureConnections` | bool | No | Allow insecure connections |

### ProviderConfig (in workflow node)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `webhookData` | string | No | Request body data (JSON string, overrides access `data`) |
| `headers` | string | No | Additional headers (merged with access headers) |
| `timeout` | int | No | Request timeout in seconds (default: 30) |

**Template Variables** (available in `webhookData`):
- `${CERTIMATE_NOTIFIER_SUBJECT}` or `${SUBJECT}` - Notification subject
- `${CERTIMATE_NOTIFIER_MESSAGE}` or `${MESSAGE}` - Notification message

**Example:**
```json
{
  "provider": "webhook",
  "providerAccessId": "ACCESS_WEBHOOK",
  "providerConfig": {
    "webhookData": "{\"title\":\"${SUBJECT}\",\"body\":\"${MESSAGE}\"}",
    "timeout": 60
  },
  "subject": "[Certimate] Alert",
  "message": "Certificate renewal completed."
}
```

---

## Template Variables

The `subject` and `message` fields support template variables using `{{$variable}}` syntax:

| Variable | Description |
|----------|-------------|
| `{{$workflow.id}}` | Workflow ID |
| `{{$workflow.name}}` | Workflow name |
| `{{$certificate.commonName}}` | Certificate common name |
| `{{$certificate.subjectAltNames}}` | Subject alternative names |
| `{{$certificate.notBefore}}` | Not before date |
| `{{$certificate.notAfter}}` | Not after date |
| `{{$certificate.daysLeft}}` | Days until expiry |
| `{{$node.skipped}}` | Whether node was skipped |

---

## Summary Table

| Provider | Access Type | Required Access Fields | ProviderConfig Fields |
|----------|-------------|------------------------|----------------------|
| `email` | email | `smtpHost`, `smtpPort`, `smtpTls`, `username`, `password`, `senderAddress` | `receiverAddress`*, `format` |
| `dingtalkbot` | dingtalkbot | `webhookUrl` | None |
| `larkbot` | larkbot | `webhookUrl` | None |
| `wecombot` | wecombot | `webhookUrl` | None |
| `slackbot` | slackbot | `botToken` | `channelId`* |
| `discordbot` | discordbot | `botToken` | `channelId`* |
| `telegrambot` | telegrambot | `botToken` | `chatId`* |
| `mattermost` | mattermost | `serverUrl`, `username`, `password` | `channelId`* |
| `webhook` | webhook | `url` | `webhookData`, `headers`, `timeout` |

*Field can be in either access config or providerConfig
