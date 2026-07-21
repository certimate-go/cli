---
name: ctm-certificate
version: 1.0.0
description: |
  Certimate CLI: Certificate management for SSL certificates.
  Use when: (1) Listing or viewing SSL certificates, (2) Checking certificate
  expiry status, (3) Downloading certificates as ZIP archives.
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

Downloads a certificate as a **ZIP archive** — the server always returns an
archive, never a single PEM/PFX file. The archive bundles the cert, private key,
and chain:

- `--format PEM` (default): `<domain>.key`, `<domain>.crt`, `<domain> (server).pem`, `<domain> (intermedia).pem`
- `--format PFX` / `--format JKS`: a single `<domain>.pfx` / `<domain>.jks`. The CLI does not expose password/alias flags, so server defaults apply (e.g. PFX password is `certimate`).

The destination is `--dest` (not `--output` — that is the global output-format flag, `json|table`, which controls the status line below):

```bash
# Default: writes <domain>-<id>.zip in the current directory
certimate certificate download CERTIFICATE_ID

# A directory (or path ending in /) writes <domain>-<id>.zip inside it
certimate certificate download CERTIFICATE_ID --dest ./certs/

# A name ending in .zip is used as-is
certimate certificate download CERTIFICATE_ID --dest cert.zip

# A non-.zip name gets .zip appended -> writes test.pem.zip
certimate certificate download CERTIFICATE_ID --dest test.pem

# Stream the raw archive bytes to stdout
certimate certificate download CERTIFICATE_ID --dest -
```

`--dest` rules:

- **Omitted** (default): writes `<domain>-<certificateId>.zip` in the current directory.
- **Directory** (existing dir, or path ending in `/`): writes `<domain>-<certificateId>.zip` inside it; the directory is created if it does not exist. `<domain>` is the first SAN with `*` replaced by `_`.
- **File path**: used verbatim if it ends in `.zip` (case-insensitive); otherwise `.zip` is appended.
- **`-`**: raw archive bytes are written to stdout.

The status output (formatted by the global `--output`/`-o`) reports the actual file written:

```json
{ "status": "downloaded", "format": "PEM", "dest": "example.com-t9vjtlj60tsbahy.zip", "size_bytes": 1234 }
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
# Download every certificate into ./backups/, each as <domain>-<id>.zip
mkdir -p backups
for id in $(certimate certificate list | jq -r '.items[].id'); do
  certimate certificate download "$id" --format PEM --dest backups/
done
```

## Certificate Status

The table output shows status:
- **valid**: Certificate is valid
- **EXPIRING**: Certificate expires within 30 days
- **EXPIRED**: Certificate has expired
