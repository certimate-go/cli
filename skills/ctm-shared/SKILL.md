---
name: ctm-shared
version: 1.0.0
description: |
  Certimate CLI: Shared configuration and common patterns.
  Use when: (1) Setting up or configuring certimate-cli, (2) Understanding
  authentication, (3) Checking version, (4) Troubleshooting connection issues.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---

# Certimate CLI - Shared Configuration

Common patterns and setup for certimate-cli.

## Installation

### From Source

```bash
go install github.com/certimate-go/cli@latest
```

### From Binary

Download from [GitHub Releases](https://github.com/certimate-go/cli/releases).

## Configuration

### Initial Setup

```bash
certimate config set --server https://certimate.example.com --token YOUR_TOKEN
```

The token is a superuser impersonation token from PocketBase admin dashboard
(Collections > `_superusers` > select user > "Impersonate").
See [PocketBase docs](https://pocketbase.io/docs/authentication/#api-keys) for details.

### Environment Variables

| Variable               | Description                  |
| ---------------------- | ---------------------------- |
| `CERTIMATE_CLI_TOKEN`  | API token (highest priority) |
| `CERTIMATE_CLI_SERVER` | Server URL                   |

### Multiple Profiles

```bash
# Create production profile
certimate config set --server https://certimate.example.com --token PROD_TOKEN --profile production

# Use profile
certimate workflow list --profile production
```

## Output Formats

```bash
# JSON (default, for AI agents)
certimate workflow list

# Table (for humans)
certimate workflow list --output table
```

## Exit Codes

| Code | Meaning              | Action                    |
| ---- | -------------------- | ------------------------- |
| 0    | Success              | Continue                  |
| 1    | General error        | Check error message       |
| 2    | Invalid arguments    | Check command syntax      |
| 3    | Authentication error | Re-run `config set`       |
| 4    | Network error        | Check server connectivity |
| 5    | Not found            | Verify resource ID        |

## Common Errors

### "no server configured"

Run: `certimate config set --server URL --token TOKEN`

### "authentication failed"

Token is invalid or expired. Generate a new token from Certimate web UI.

### "connection refused"

Server is not reachable. Check URL and network connectivity.

## Version Check

```bash
certimate version
```
