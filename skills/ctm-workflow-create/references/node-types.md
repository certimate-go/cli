# Node Types Reference

This document describes all available node types and their configurations, based on the actual source code at `~/work/certimate/internal/domain/workflow.go`.

> **Important**: The JSON field names match the Go struct tags exactly. The provider type field is `provider` (not `providerType`), and the provider-specific config field is `providerConfig` (not `extendedConfig`).

## Core Nodes

### start

Entry point for the workflow.

```json
{
  "id": "start-node",
  "type": "start",
  "data": {
    "name": "Start",
    "config": {
      "trigger": "manual"
    }
  }
}
```

**Config fields**:
- `trigger` (string, required): `"manual"` or `"scheduled"`

### end

Termination point for the workflow.

```json
{
  "id": "end-node",
  "type": "end",
  "data": {
    "name": "End"
  }
}
```

No config fields required.

## Control Flow Nodes

### tryCatch

Wraps a try/catch block for error handling.

```json
{
  "id": "try-catch-node",
  "type": "tryCatch",
  "data": { "name": "Try Execute" },
  "blocks": [
    { "type": "tryBlock", ... },
    { "type": "catchBlock", ... }
  ]
}
```

**Structure**: Must contain exactly one `tryBlock` and one `catchBlock` in `blocks` array.

### tryBlock

Contains the main execution path.

```json
{
  "id": "try-block",
  "type": "tryBlock",
  "data": { "name": "" },
  "blocks": [
    { "type": "bizApply", ... },
    { "type": "bizDeploy", ... }
  ]
}
```

**Structure**: Contains business logic nodes in `blocks` array.

### catchBlock

Contains error handling logic.

```json
{
  "id": "catch-block",
  "type": "catchBlock",
  "data": { "name": "On Error" },
  "blocks": [
    { "type": "bizNotify", ... },
    { "type": "end", ... }
  ]
}
```

**Structure**: Contains error handling nodes in `blocks` array.

### delay

Wait for a specified number of seconds.

```json
{
  "id": "delay-node",
  "type": "delay",
  "data": {
    "name": "Wait",
    "config": {
      "wait": 60
    }
  }
}
```

**Config fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `wait` | int | Yes | Wait time in seconds |

### condition

Conditional node. No specific config struct.

```json
{
  "id": "condition-node",
  "type": "condition",
  "data": { "name": "Condition" }
}
```

### branchBlock

Branch with an expression condition.

```json
{
  "id": "branch-node",
  "type": "branchBlock",
  "data": {
    "name": "Branch",
    "config": {
      "expression": "..."
    }
  }
}
```

**Config fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `expression` | string | Yes | Condition expression |

## Business Nodes

### bizApply (Certificate Application)

Applies for an SSL certificate using ACME protocol. This is the core node for certificate issuance.

```json
{
  "id": "apply-node",
  "type": "bizApply",
  "data": {
    "name": "Apply Certificate",
    "config": {
      "identifier": "domain",
      "domains": "example.com;*.example.com",
      "ipaddrs": "",
      "contactEmail": "admin@example.com",
      "challengeType": "dns-01",
      "provider": "aliyun-dns",
      "providerAccessId": "ACCESS_ID",
      "providerConfig": {},
      "caProvider": "",
      "caProviderAccessId": "",
      "caProviderConfig": {},
      "keySource": "auto",
      "keyAlgorithm": "RSA2048",
      "keyContent": "",
      "validityLifetime": "",
      "preferredChain": "",
      "acmeProfile": "",
      "nameservers": "",
      "dnsPropagationWait": 0,
      "dnsPropagationTimeout": 0,
      "dnsTTL": 0,
      "httpDelayWait": 0,
      "disableCommonName": false,
      "disableFollowCNAME": false,
      "disableARI": false,
      "skipBeforeExpiryDays": 30
    }
  }
}
```

**Required fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `identifier` | string | Yes | `"domain"` or `"ip"` - determines whether to use domains or ipaddrs |
| `domains` | string | Conditional | Domains to include (semicolon-separated). Required if identifier="domain" |
| `ipaddrs` | string | Conditional | IP addresses to include (semicolon-separated). Required if identifier="ip" |
| `contactEmail` | string | Yes | ACME contact email address |
| `challengeType` | string | Yes | Challenge type: `dns-01` or `http-01` |
| `provider` | string | Yes | ACME challenge provider (see below) |
| `providerAccessId` | string | Yes | Access credential ID for the provider |
| `keySource` | string | Yes | Key source: `auto`, `reuse`, `custom` |
| `keyAlgorithm` | string | Yes | Key algorithm |

**Optional fields**:
| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `providerConfig` | object | `{}` | Provider extra config |
| `caProvider` | string | `""` | CA provider (empty = global config) |
| `caProviderAccessId` | string | `""` | CA provider access record ID |
| `caProviderConfig` | object | `{}` | CA provider extra config |
| `keyContent` | string | `""` | PEM private key (only if `keySource: "custom"`) |
| `validityLifetime` | string | `""` | Validity period, e.g. `"30d"`, `"6h"` |
| `preferredChain` | string | `""` | Preferred certificate chain |
| `acmeProfile` | string | `""` | ACME Profiles Extension |
| `nameservers` | string | `""` | Custom DNS servers (semicolon-separated) |
| `dnsPropagationWait` | int | `0` | DNS propagation wait (seconds) |
| `dnsPropagationTimeout` | int | `0` | DNS propagation check timeout |
| `dnsTTL` | int | `0` | DNS record TTL |
| `httpDelayWait` | int | `0` | HTTP challenge delay wait |
| `disableCommonName` | bool | `false` | Exclude CommonName from certificate |
| `disableFollowCNAME` | bool | `false` | Disable CNAME following |
| `disableARI` | bool | `false` | Disable ARI (ACME Renewal Information) |
| `skipBeforeExpiryDays` | int | `30` | Days before expiry to skip renewal |

**Challenge types**:
- `dns-01` — DNS-based challenge (recommended for wildcards)
- `http-01` — HTTP-based challenge (requires web server)

**Key source values**:
- `auto` — Generate new private key automatically (default)
- `reuse` — Reuse private key from previous certificate
- `custom` — Use a user-provided private key (set `keyContent`)

**Key algorithms**:
| Algorithm | Description |
|-----------|-------------|
| `RSA2048` | RSA 2048-bit (most compatible, default) |
| `RSA3072` | RSA 3072-bit |
| `RSA4096` | RSA 4096-bit |
| `RSA8192` | RSA 8192-bit |
| `EC256` | ECDSA P-256 (modern, smaller keys) |
| `EC384` | ECDSA P-384 |
| `EC512` | ECDSA P-521 |

**DNS-01 providers** (75+ providers, `provider` field values):
`35cn`, `51dnscom`, `acmedns`, `acmehttpreq`, `akamai-edgedns`, `aliyun-dns`, `aliyun-esa`, `arvancloud`, `aws-route53`, `azure-dns`, `baiducloud-dns`, `cloudflare`, `cloudns`, `desec`, `digitalocean`, `godaddy`, `hetzner`, `huaweicloud-dns`, `jdcloud-dns`, `namedotcom`, `namesilo`, `netcup`, `ns1`, `ovhcloud`, `porkbun`, `powerdns`, `qingcloud-dns`, `rfc2136`, `spaceship`, `tencentcloud-dns`, `tencentcloud-eo`, `volcengine-dns`, `vultr`, `westcn`, `xinnet`, ...

**HTTP-01 providers**:
- `local` — Run HTTP server locally
- `s3` — Upload challenge to S3-compatible storage
- `ssh` — Upload challenge via SSH

**CA providers** (`caProvider` field):
| Value | Description |
|-------|-------------|
| `""` (empty) | Use global config |
| `letsencrypt` | Let's Encrypt (default) |
| `letsencryptstaging` | Let's Encrypt Staging |
| `zerossl` | ZeroSSL |
| `googletrustservices` | Google Trust Services |
| `digicert` | DigiCert |
| `sectigo` | Sectigo |
| `sslcom` | SSL.com |
| `acmeca` | ACME CA |
| `actalisssl` | Actalis |
| `globalsignatlas` | GlobalSign Atlas |
| `litessl` | LiteSSL |

**Node outputs** (available as template variables in downstream nodes):
- `certificate.commonName` — Common name
- `certificate.subjectAltNames` — Subject alternative names
- `certificate.notBefore` — Not before date
- `certificate.notAfter` — Not after date
- `certificate.hoursLeft` — Hours until expiry
- `certificate.daysLeft` — Days until expiry
- `certificate.validity` — Validity status
- `node.skipped` — Whether the node was skipped

### bizUpload (Certificate Upload)

Upload/import an existing certificate (instead of applying via ACME).

```json
{
  "id": "upload-node",
  "type": "bizUpload",
  "data": {
    "name": "Upload Certificate",
    "config": {
      "source": "form",
      "certificate": "-----BEGIN CERTIFICATE-----\n...",
      "privateKey": "-----BEGIN PRIVATE KEY-----\n..."
    }
  }
}
```

**Config fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `source` | string | Yes | Source type: `form`, `local`, `url` |
| `certificate` | string | Yes | PEM content, file path, or URL |
| `privateKey` | string | Yes | PEM content, file path, or URL |

**Source values**:
- `form` — PEM content entered directly in `certificate` and `privateKey` fields
- `local` — File path on local filesystem
- `url` — Downloadable URL

**Node outputs**: Same certificate variables as bizApply.

### bizMonitor (Certificate Monitoring)

Monitor a live certificate's expiry status.

```json
{
  "id": "monitor-node",
  "type": "bizMonitor",
  "data": {
    "name": "Monitor Certificate",
    "config": {
      "host": "example.com",
      "port": 443,
      "domain": "",
      "requestPath": ""
    }
  }
}
```

**Config fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Host address |
| `port` | int32 | No | Port (default: 443) |
| `domain` | string | No | SNI domain (default: host) |
| `requestPath` | string | No | Request path |

**Node outputs**: Same certificate variables as bizApply (`certificate.commonName`, etc.).

### bizDeploy (Certificate Deployment)

Deploys a certificate to a target service.

```json
{
  "id": "deploy-node",
  "type": "bizDeploy",
  "data": {
    "name": "Deploy Certificate",
    "config": {
      "provider": "aliyun-cdn",
      "providerAccessId": "ACCESS_ID",
      "providerConfig": {
        "region": "cn-hangzhou",
        "domain": "cdn.example.com",
        "domainMatchPattern": "exact"
      },
      "skipOnLastSucceeded": true,
      "certificateOutputNodeId": "apply-node"
    }
  }
}
```

**Config fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `provider` | string | Yes | Deployment provider type (see below) |
| `providerAccessId` | string | Yes | Access credential ID |
| `certificateOutputNodeId` | string | Yes | ID of upstream bizApply/bizUpload node |
| `providerConfig` | object | Yes | Provider-specific configuration |
| `skipOnLastSucceeded` | bool | No | Skip if last deployment succeeded (default: `true`) |

**120+ deploy providers**: See `deploy-configs.md` for the full list and their `providerConfig` fields.

### bizNotify (Notification)

Sends a notification via configured channel.

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
      "message": "Workflow completed.",
      "skipOnAllPrevSkipped": false
    }
  }
}
```

**Config fields**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `provider` | string | Yes | Notification provider type (see below) |
| `providerAccessId` | string | Yes | Access credential ID |
| `subject` | string | Yes | Notification subject |
| `message` | string | Yes | Notification body (supports `{{$variable}}` templates) |
| `providerConfig` | object | No | Provider-specific extended config |
| `skipOnAllPrevSkipped` | bool | No | Skip if all previous nodes were skipped |

**Notification providers**:
| Provider | Access Type | ProviderConfig Fields |
|----------|-------------|----------------------|
| `email` | email | `receiverAddress`*, `format` |
| `dingtalkbot` | dingtalkbot | None |
| `larkbot` | larkbot | None |
| `wecombot` | wecombot | None |
| `slackbot` | slackbot | `channelId`* |
| `discordbot` | discordbot | `channelId`* |
| `telegrambot` | telegrambot | `chatId`* |
| `mattermost` | mattermost | `channelId`* |
| `webhook` | webhook | `webhookData`, `headers`, `timeout` |

*Field can be in either access config or providerConfig

**See `notify-configs.md` for detailed configuration of each provider.**

**Template variables** (in `subject` and `message`):
- `{{$workflow.id}}` - Workflow ID
- `{{$workflow.name}}` - Workflow name
- `{{$certificate.commonName}}` - Certificate common name
- `{{$certificate.subjectAltNames}}` - Subject alternative names
- `{{$certificate.notBefore}}` - Not before date
- `{{$certificate.notAfter}}` - Not after date
- `{{$certificate.daysLeft}}` - Days until expiry

## Node ID Best Practices

1. Use descriptive names: `apply-cert`, `deploy-to-cdn`, `notify-slack`
2. Keep IDs unique within a workflow
3. Reference IDs correctly in `certificateOutputNodeId`
4. Use consistent naming conventions

## Common Patterns

### Pattern 1: Simple Apply Only (No Error Handling)

```json
{
  "nodes": [
    { "id": "start", "type": "start", "data": { "name": "Start", "config": { "trigger": "manual" } } },
    { "id": "apply", "type": "bizApply", "data": { "name": "Apply", "config": { ... } } },
    { "id": "end", "type": "end", "data": { "name": "End" } }
  ]
}
```

### Pattern 2: Apply and Deploy (No Error Handling)

```json
{
  "nodes": [
    { "id": "start", "type": "start", "data": { "name": "Start", "config": { "trigger": "manual" } } },
    { "id": "apply", "type": "bizApply", "data": { "name": "Apply", "config": { ... } } },
    { "id": "deploy", "type": "bizDeploy", "data": { "name": "Deploy", "config": { "certificateOutputNodeId": "apply", ... } } },
    { "id": "end", "type": "end", "data": { "name": "End" } }
  ]
}
```

### Pattern 3: Full Workflow with Error Handling

```json
{
  "nodes": [
    { "id": "start", "type": "start", "data": { "name": "Start", "config": { "trigger": "manual" } } },
    {
      "id": "trycatch",
      "type": "tryCatch",
      "data": { "name": "Try" },
      "blocks": [
        {
          "id": "try",
          "type": "tryBlock",
          "data": { "name": "" },
          "blocks": [
            { "id": "apply", "type": "bizApply", "data": { "name": "Apply", "config": { ... } } },
            { "id": "deploy", "type": "bizDeploy", "data": { "name": "Deploy", "config": { ... } } }
          ]
        },
        {
          "id": "catch",
          "type": "catchBlock",
          "data": { "name": "On Error" },
          "blocks": [
            { "id": "notify", "type": "bizNotify", "data": { "name": "Notify", "config": { ... } } },
            { "id": "catch-end", "type": "end", "data": { "name": "End" } }
          ]
        }
      ]
    },
    { "id": "end", "type": "end", "data": { "name": "End" } }
  ]
}
```

### Pattern 4: Multiple Deploys

Multiple `bizDeploy` nodes in sequence within the same `tryBlock.blocks` array:

```json
{
  "blocks": [
    { "id": "apply", "type": "bizApply", ... },
    { "id": "deploy1", "type": "bizDeploy", "data": { "config": { "certificateOutputNodeId": "apply", "provider": "aliyun-cdn", ... } } },
    { "id": "deploy2", "type": "bizDeploy", "data": { "config": { "certificateOutputNodeId": "apply", "provider": "ssh", ... } } }
  ]
}
```

## Complete Workflow Example

**IMPORTANT**: There is NO `edges` array. Connections are implicit through the nested `blocks` structure. Sequential execution is determined by array order.

```json
{
  "nodes": [
    {
      "id": "start-node",
      "type": "start",
      "data": {
        "name": "Start",
        "config": { "trigger": "manual" }
      }
    },
    {
      "id": "try-catch-node",
      "type": "tryCatch",
      "data": { "name": "Try Execute" },
      "blocks": [
        {
          "id": "try-block",
          "type": "tryBlock",
          "data": { "name": "" },
          "blocks": [
            {
              "id": "apply-node",
              "type": "bizApply",
              "data": {
                "name": "Apply Certificate",
                "config": {
                  "domains": ["example.com", "*.example.com"],
                  "challengeType": "dns-01",
                  "provider": "aliyun-dns",
                  "providerAccessId": "ACCESS_ALIYUN_DNS",
                  "keySource": "auto",
                  "keyAlgorithm": "RSA2048",
                  "skipBeforeExpiryDays": 30
                }
              }
            },
            {
              "id": "deploy-node",
              "type": "bizDeploy",
              "data": {
                "name": "Deploy to CDN",
                "config": {
                  "provider": "aliyun-cdn",
                  "providerAccessId": "ACCESS_ALIYUN_CDN",
                  "skipOnLastSucceeded": true,
                  "certificateOutputNodeId": "apply-node",
                  "providerConfig": {
                    "region": "cn-hangzhou",
                    "domain": "cdn.example.com",
                    "domainMatchPattern": "exact"
                  }
                }
              }
            }
          ]
        },
        {
          "id": "catch-block",
          "type": "catchBlock",
          "data": { "name": "On Error" },
          "blocks": [
            {
              "id": "notify-node",
              "type": "bizNotify",
              "data": {
                "name": "Notify",
                "config": {
                  "provider": "email",
                  "providerAccessId": "ACCESS_EMAIL",
                  "providerConfig": {
                    "receiverAddress": "admin@example.com"
                  },
                  "subject": "[Certimate] Workflow Failed",
                  "message": "Workflow {{$workflow.name}} failed."
                }
              }
            },
            {
              "id": "catch-end-node",
              "type": "end",
              "data": { "name": "End" }
            }
          ]
        }
      ]
    },
    {
      "id": "end-node",
      "type": "end",
      "data": { "name": "End" }
    }
  ]
}
```
