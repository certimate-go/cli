---
name: ctm-certificate
version: 1.0.0
description: |
  Certimate CLI: Certificate management for SSL certificates.
  Use when: (1) Listing or viewing SSL certificates, (2) Checking certificate
  expiry status, (3) Downloading certificates in various formats.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---
# Certimate Certificate Management

Manage SSL certificates through the command line.

## Available Commands

### List Certificates

```bash
certimate certificate list [--output json|table] [--filter EXPR]
```

**Example Output (JSON):**
```json
{
  "items": [
    {
      "id": "cert_abc123",
      "subjectAltNames": ["example.com", "*.example.com"],
      "issuer": "Let's Encrypt",
      "notBefore": "2026-01-01T00:00:00Z",
      "notAfter": "2026-04-01T00:00:00Z",
      "workflowId": "wf_xyz"
    }
  ]
}
```

**Table Output:**
```
ID           SUBJECT                    EXPIRES     STATUS    CREATED
cert_abc123  example.com, *.example.com 2026-04-01  valid     2026-01-01
```

### Get Certificate Details

```bash
certimate certificate get CERTIFICATE_ID
```

### Download Certificate

```bash
# Download as PEM
certimate certificate download CERTIFICATE_ID --format PEM

# Download to file
certimate certificate download CERTIFICATE_ID --format PEM --output cert.pem

# Download as PFX
certimate certificate download CERTIFICATE_ID --format PFX --output cert.pfx
```

## Filtering

Use PocketBase filter syntax:

```bash
# Certificates expiring before a date
certimate certificate list --filter 'notAfter<"2026-04-01"'

# Certificates from specific workflow
certimate certificate list --filter 'workflowId="wf_xyz"'
```

## Common Patterns

### Check Expiring Certificates

```bash
# List all certificates with table output
certimate certificate list --output table

# Find certificates expiring within 30 days
certimate certificate list | jq '.items[] | select(.notAfter < (now + 2592000 | strftime("%Y-%m-%dT%H:%M:%SZ")))'
```

### Renew Certificate

```bash
# 1. Find the certificate
certimate certificate get CERTIFICATE_ID

# 2. Note the workflowId from the output

# 3. Run the workflow to renew
certimate workflow run WORKFLOW_ID --wait
```

### Backup Certificates

```bash
# Download all certificates
for id in $(certimate certificate list | jq -r '.items[].id'); do
  certimate certificate download "$id" --format PEM --output "cert_$id.pem"
done
```

## Certificate Status

The table output shows status:
- **valid**: Certificate is valid
- **EXPIRING**: Certificate expires within 30 days
- **EXPIRED**: Certificate has expired
