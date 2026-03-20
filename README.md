# certimate-cli

Command-line interface for [Certimate](https://github.com/certimate-go/certimate) SSL certificate management.

[![Go Report Card](https://goreportcard.com/badge/github.com/certimate-go/cli)](https://goreportcard.com/report/github.com/certimate-go/cli)
[![License](https://img.shields.io/github/license/certimate-go/cli)](LICENSE)

## Features

- Workflow management: list, view, and execute certificate workflows
- Certificate management: view and download SSL certificates in multiple formats
- Access management: create, edit, and delete provider credentials
- Multiple output formats: JSON (default) and table output
- Profile support: manage multiple server configurations
- AI agent ready: includes SKILL.md files for AI agent integration

## Installation

### From Source

```bash
go install github.com/certimate-go/cli@latest
```

### From Binary

Download the latest release from [GitHub Releases](https://github.com/certimate-go/cli/releases).

### Using Homebrew

```bash
brew install certimate-go/tap/certimate
```

## Quick Start

### 1. Configure Connection

```bash
certimate config set --server http://127.0.0.1:8090 --token YOUR_API_TOKEN
```

Obtain an API token from Certimate web UI (Settings > API Tokens).

### 2. List Workflows

```bash
certimate workflow list
```

### 3. Execute a Workflow

```bash
# Fire-and-forget
certimate workflow run WORKFLOW_ID

# Wait for completion
certimate workflow run WORKFLOW_ID --wait
```

### 4. List Certificates

```bash
certimate certificate list

# With table output
certimate certificate list --output table
```

## Commands

### Configuration

```bash
# Set configuration
certimate config set --server URL --token TOKEN [--profile NAME]

# View current configuration
certimate config get

# Show current profile
certimate config current

# List all profiles
certimate config list
```

### Workflows

```bash
# List workflows
certimate workflow list [--filter EXPR] [--limit N]

# Get workflow details
certimate workflow get WORKFLOW_ID

# Execute workflow
certimate workflow run WORKFLOW_ID [--wait] [--timeout SECONDS]

# Cancel running workflow
certimate workflow cancel WORKFLOW_ID RUN_ID

# List execution history
certimate workflow runs WORKFLOW_ID [--limit N]
```

### Certificates

```bash
# List certificates
certimate certificate list [--filter EXPR] [--limit N]

# Get certificate details
certimate certificate get CERTIFICATE_ID

# Download certificate
certimate certificate download CERTIFICATE_ID --format PEM|PFX|JKS [--output FILE]
```

### Access Credentials

```bash
# List access credentials
certimate access list [--reveal]

# Get access details
certimate access get ACCESS_ID [--reveal]

# Create new access credential
certimate access create --name NAME --provider PROVIDER --config JSON

# Edit access credential
certimate access edit ACCESS_ID --name NAME --config JSON

# Delete access credential
certimate access delete ACCESS_ID
```

### Other Commands

```bash
# Show version
certimate version

# Generate shell completion
certimate completion bash|zsh|fish
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Config file path (default `~/.config/certimate-cli/config.yaml`) |
| `--debug` | Enable debug output |
| `-o, --output` | Output format: `json` or `table` (default `json`) |
| `--profile` | Configuration profile (default `default`) |

## Output Formats

All commands support JSON (default) and table output:

```bash
# JSON output (default)
certimate workflow list

# Table output
certimate workflow list --output table
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CERTIMATE_CLI_TOKEN` | API token (overrides config file) |
| `CERTIMATE_CLI_SERVER` | Server URL (overrides config file) |
| `CERTIMATE_CLI_CONFIG_DIR` | Custom config directory |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Authentication error |
| 4 | Network error |
| 5 | Not found |

## AI Agent Skills

The repo ships Agent Skills (SKILL.md files) for AI-powered SSL certificate management:

- `ctm-shared` - Common patterns, authentication, and CLI setup
- `ctm-workflow` - Certificate workflow operations
- `ctm-certificate` - Certificate viewing and downloading
- `ctm-access` - Provider credential management

### Using npx

```bash
# Install all skills at once
npx skills add https://github.com/certimate-go/cli

# Or pick only what you need
npx skills add https://github.com/certimate-go/cli/tree/main/skills/ctm-workflow
npx skills add https://github.com/certimate-go/cli/tree/main/skills/ctm-certificate
```

### Manual Installation

```bash
cp -r skills/* ~/.claude/skills/
```

## Development

```bash
# Build
make build

# Install locally
make install

# Run tests
make test

# Build for all platforms
make build-all

# Generate shell completion
make completion
```

## License

[MIT License](LICENSE)

## Related Projects

- [Certimate](https://github.com/certimate-go/certimate) - SSL certificate management system
