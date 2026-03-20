---
name: ctm-access
version: 2.0.0
description: |
  Certimate CLI: Access credentials management for DNS providers, CDNs,
  and other services. Use when: (1) Listing provider credentials,
  (2) Viewing access configuration, (3) Creating new access credentials,
  (4) Editing existing access credentials, (5) Deleting access credentials.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---
# Certimate Access Management

Manage provider access credentials through the command line.

## Security Note

Access credentials contain sensitive information (API keys, secrets, passwords).
By default, sensitive fields are masked as `********`. Use `--reveal` flag to show actual values.

## Step-by-Step Access Creation

The skill supports interactive step-by-step access creation:

### Step 1: Provide Basic Information

First, provide the access name and provider type:

```bash
certimate access create --name "My Cloudflare Access" --provider "cloudflare"
```

### Step 2: Provide Provider-Specific Configuration

After specifying name and provider, provide the configuration. The configuration format depends on the provider:

**Cloudflare Example:**
```bash
certimate access create \
  --name "My Cloudflare Access" \
  --provider "cloudflare" \
  --config '{"dnsApiToken": "your-dns-api-token"}'
```

**Aliyun Example:**
```bash
certimate access create \
  --name "My Aliyun Access" \
  --provider "aliyun" \
  --config '{"accessKeyId": "your-access-key-id", "accessKeySecret": "your-access-key-secret"}'
```

**AWS Route53 Example:**
```bash
certimate access create \
  --name "My AWS Access" \
  --provider "aws" \
  --config '{"accessKeyId": "AKIAIOSFODNN7EXAMPLE", "secretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}'
```

### Step 3: Complete Creation

The command will create the access and return the result (with sensitive fields masked by default).

## Available Commands

### List Access Credentials

```bash
# List with masked secrets (default)
certimate access list [--output json|table]

# Show actual values (use with caution)
certimate access list --reveal
```

**Example Output (masked):**
```json
{
  "items": [
    {
      "id": "access_abc123",
      "name": "Cloudflare Production",
      "provider": "cloudflare",
      "config": {
        "dnsApiToken": "********"
      },
      "created": "2026-01-01T00:00:00Z"
    }
  ]
}
```

### Get Access Details

```bash
# Masked (default)
certimate access get ACCESS_ID

# Show actual values
certimate access get ACCESS_ID --reveal
```

### Create Access

```bash
# From JSON string
certimate access create \
  --name "Cloudflare Production" \
  --provider "cloudflare" \
  --config '{"dnsApiToken": "secret-token"}'

# From JSON file
certimate access create \
  --name "Cloudflare Production" \
  --provider "cloudflare" \
  --config @/path/to/config.json
```

**Example Output (masked):**
```json
{
  "id": "access_xyz789",
  "name": "Cloudflare Production",
  "provider": "cloudflare",
  "config": {
    "dnsApiToken": "********"
  },
  "created": "2026-03-19T10:00:00Z",
  "updated": "2026-03-19T10:00:00Z"
}
```

### Edit Access

```bash
# Update name only
certimate access edit ACCESS_ID --name "Cloudflare Staging"

# Update config
certimate access edit ACCESS_ID --config '{"dnsApiToken": "new-secret-token"}'

# Update multiple fields
certimate access edit ACCESS_ID \
  --name "Cloudflare Staging" \
  --provider "cloudflare" \
  --config '{"dnsApiToken": "new-secret-token"}'
```

### Delete Access

```bash
certimate access delete ACCESS_ID
```

**Example Output:**
```json
{
  "status": "deleted",
  "id": "access_xyz789",
  "message": "Access credential deleted successfully"
}
```

## Provider Configuration Reference

For detailed configuration requirements for each provider, see the [Provider Configuration Reference](references/providers/README.md).

### Common Provider Examples

**Cloudflare DNS API Token:**
```json
{
  "dnsApiToken": "your-cloudflare-dns-api-token"
}
```

**Aliyun Access Key:**
```json
{
  "accessKeyId": "your-access-key-id",
  "accessKeySecret": "your-access-key-secret"
}
```

**AWS Route53:**
```json
{
  "accessKeyId": "AKIAIOSFODNN7EXAMPLE",
  "secretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

**Tencent Cloud:**
```json
{
  "secretId": "your-secret-id",
  "secretKey": "your-secret-key"
}
```

## Supported Providers

Certimate supports 100+ providers including:

### DNS Providers
- Cloudflare
- Aliyun (Alibaba Cloud)
- Tencent Cloud
- AWS Route53
- GoDaddy
- Namecheap
- ... and more

### CDN/Deployment Targets
- Aliyun CDN
- Tencent Cloud CDN
- AWS CloudFront
- Kubernetes
- SSH servers
- ... and more

### Notification Providers
- Email
- DingTalk Bot
- Slack Bot
- Telegram Bot
- Webhook
- ... and more

## Common Patterns

### List All Cloudflare Accesses

```bash
certimate access list --filter 'provider="cloudflare"'
```

### Find Access by Name

```bash
certimate access list | jq '.items[] | select(.name | contains("Production"))'
```

### Audit Access Credentials

```bash
# List all accesses with table output
certimate access list --output table

# Check which accesses are being used by workflows
# (requires cross-referencing with workflow configs)
```

## Security Best Practices

1. **Never log revealed credentials**: Avoid piping `--reveal` output to logs
2. **Use environment variables**: Set `CERTIMATE_CLI_TOKEN` securely
3. **Rotate credentials regularly**: Update access credentials periodically
4. **Principle of least privilege**: Ensure each access has minimal required permissions
5. **Use file-based config for sensitive data**: Store config in files rather than command line arguments

## Configuration File Format

Access configuration can be provided as:
1. **JSON string**: `--config '{"key": "value"}'`
2. **JSON file**: `--config @/path/to/config.json`

### Example Config File (cloudflare.json)
```json
{
  "dnsApiToken": "your-cloudflare-dns-api-token"
}
```

### Usage with File
```bash
certimate access create \
  --name "Cloudflare Production" \
  --provider "cloudflare" \
  --config @cloudflare.json
```

## Troubleshooting

### Common Issues

1. **Invalid JSON format**: Ensure your config is valid JSON
2. **Missing required fields**: Check provider-specific requirements in references/providers/README.md
3. **Authentication errors**: Verify your API credentials are correct and have required permissions
4. **Network errors**: Ensure Certimate server is accessible

### Getting Help

- View command help: `certimate access --help`
- View specific command help: `certimate access create --help`
- Check provider config: See references/providers/README.md
