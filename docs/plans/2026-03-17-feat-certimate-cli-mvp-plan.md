---
title: certimate-cli MVP Implementation
type: feat
date: 2026-03-17
---

# certimate-cli MVP Implementation

## Overview

Build a command-line interface (CLI) tool for Certimate that enables certificate workflow management, certificate viewing, and authorization management through the terminal. The CLI follows the googleworkspace/cli architecture pattern with AI agent skills support via SKILL.md files.

## Problem Statement / Motivation

Certimate is a self-hosted SSL certificate management system with a web UI, but lacks command-line access. This creates friction for:

1. **DevOps engineers** who need to integrate certificate management into CI/CD pipelines
2. **AI agents** that need programmatic access to certificate workflows
3. **Power users** who prefer terminal-based workflows

The CLI bridges this gap by providing a complete command surface with structured JSON output for automation.

## Proposed Solution

Build a Go-based CLI using `spf13/cobra` that:
1. Connects to Certimate server via HTTP REST API (PocketBase)
2. Provides commands for workflow, certificate, and access management
3. Includes SKILL.md files for AI agent integration
4. Supports JSON (default) and table output formats

### Architecture

```
certimate-cli/
├── cmd/
│   ├── root.go              # Root command, global flags
│   ├── config.go            # Configuration management
│   ├── workflow.go          # Workflow commands
│   ├── certificate.go       # Certificate commands
│   ├── access.go            # Access commands
│   └── version.go           # Version info
├── internal/
│   ├── api/
│   │   ├── client.go        # HTTP client for Certimate API
│   │   ├── pocketbase.go    # PocketBase CRUD API wrapper
│   │   └── errors.go        # API error handling
│   ├── config/
│   │   ├── config.go        # Viper configuration
│   │   └── auth.go          # Authentication handling
│   ├── models/
│   │   ├── workflow.go      # Workflow model
│   │   ├── certificate.go   # Certificate model
│   │   ├── access.go        # Access model
│   │   └── workflow_run.go  # WorkflowRun model
│   └── output/
│       ├── json.go          # JSON formatter
│       ├── table.go         # Table formatter
│       └── output.go        # Output interface
├── skills/
│   ├── ctm-shared/
│   │   └── SKILL.md         # Common patterns, install
│   ├── ctm-workflow/
│   │   └── SKILL.md         # Workflow management
│   ├── ctm-certificate/
│   │   └── SKILL.md         # Certificate operations
│   └── ctm-access/
│       └── SKILL.md         # Access management
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### API Integration Strategy

The CLI integrates with two API surfaces:

| API | Purpose | Endpoint Pattern |
|-----|---------|------------------|
| **PocketBase CRUD** | List/Get operations | `/api/collections/{collection}/records` |
| **Custom REST API** | Actions (run, cancel, download) | `/api/workflows/{id}/runs`, `/api/certificates/{id}/download` |

**PocketBase Query Parameters:**
- `?page=1&perPage=30` - Pagination
- `?filter=(field='value')` - Filtering
- `?sort=-created` - Sorting
- `?fields=id,name` - Field selection

## Technical Considerations

### Authentication

**Token Acquisition:** Users obtain API tokens from Certimate web UI (Settings > API Tokens) as superuser tokens.

**Storage Locations (Priority):**
1. `--token` flag (lowest priority, visible in process list)
2. `CERTIMATE_CLI_TOKEN` environment variable
3. `CERTIMATE_CLI_CREDENTIALS_FILE` environment variable
4. Config file: `~/.config/certimate-cli/config.yaml`

**Config File Format (YAML):**
```yaml
current_profile: default
profiles:
  default:
    server: http://127.0.0.1:8090
    token: <encrypted-or-plain-token>
  production:
    server: https://certimate.example.com
    token: <encrypted-or-plain-token>
```

### Security

1. **Token Storage:** Store tokens in config file (consider keychain integration for v2)
2. **Access Config Masking:** `access list` masks sensitive fields by default, use `--reveal` flag
3. **Self-Signed Certs:** Support `--insecure` flag for internal deployments

### Output Format

**JSON (Default):**
```json
{
  "id": "abc123",
  "name": "Production Certificate",
  "trigger": "manual",
  "enabled": true
}
```

**Table (with `--output table`):**
```
ID       NAME                     TRIGGER   ENABLED
abc123   Production Certificate   manual    true
def456   Staging Certificate      scheduled false
```

### Error Handling

**Exit Codes:**
- `0` - Success
- `1` - General error
- `2` - Invalid arguments
- `3` - Authentication error
- `4` - Network error
- `5` - Not found

**Error Format (JSON):**
```json
{
  "error": "authentication_failed",
  "message": "Invalid or expired API token",
  "details": {
    "status_code": 401
  }
}
```

## Acceptance Criteria

### Phase 1: Project Setup

- [x] Initialize Go module `github.com/certimate-go/cli`
- [x] Create cobra CLI structure with root command
- [x] Implement Viper configuration management
- [x] Create Makefile with build/install targets
- [x] Add `.goreleaser.yml` for releases

### Phase 2: Configuration & Auth

- [x] `certimate config set --server URL --token TOKEN` - Store configuration
- [x] `certimate config get` - View current configuration
- [x] `certimate config current` - Show active profile
- [x] Support `CERTIMATE_CLI_TOKEN` and `CERTIMATE_CLI_SERVER` env vars
- [x] Validate token by making test API call on config set

### Phase 3: Workflow Commands

- [x] `certimate workflow list [--output json|table]` - List all workflows
- [x] `certimate workflow get WORKFLOW_ID` - Get workflow details
- [x] `certimate workflow run WORKFLOW_ID [--wait]` - Execute workflow
- [x] `certimate workflow cancel WORKFLOW_ID RUN_ID` - Cancel running workflow
- [x] `certimate workflow runs WORKFLOW_ID` - List workflow execution history

### Phase 4: Certificate Commands

- [x] `certimate certificate list [--output json|table]` - List certificates
- [x] `certimate certificate get CERTIFICATE_ID` - Get certificate details
- [x] `certimate certificate download CERTIFICATE_ID --format PEM|PFX|JKS` - Download certificate
- [x] Display expiry status in certificate list output

### Phase 5: Access Commands

- [x] `certimate access list [--output json|table]` - List access credentials
- [x] Mask sensitive fields (API keys, secrets) by default
- [x] Add `--reveal` flag to show full configuration

### Phase 6: AI Agent Skills

- [x] Create `skills/ctm-shared/SKILL.md` - Common patterns
- [x] Create `skills/ctm-workflow/SKILL.md` - Workflow management
- [x] Create `skills/ctm-certificate/SKILL.md` - Certificate operations
- [x] Create `skills/ctm-access/SKILL.md` - Access management
- [x] Include JSON output examples in each skill

### Phase 7: Polish

- [x] Add `certimate version` command
- [x] Add shell completion support (bash/zsh/fish)
- [x] Add `--debug` flag for verbose output
- [x] Write README.md with installation and usage
- [ ] Test with Certimate v0.3.x+

## MVP

### cmd/root.go

```go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile   string
    outputFmt string
    profile   string
)

var rootCmd = &cobra.Command{
    Use:   "certimate",
    Short: "CLI for Certimate SSL certificate management",
    Long: `Certimate CLI manages SSL certificate workflows, certificates,
and access credentials through the command line.

Configure with: certimate config set --server URL --token TOKEN
Then use: certimate workflow list, certimate certificate list, etc.`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/certimate-cli/config.yaml)")
    rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "json", "output format (json|table)")
    rootCmd.PersistentFlags().StringVar(&profile, "profile", "default", "configuration profile")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        viper.AddConfigPath("$HOME/.config/certimate-cli")
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
    }

    viper.SetEnvPrefix("CERTIMATE_CLI")
    viper.AutomaticEnv()

    viper.ReadInConfig()
}
```

### internal/api/client.go

```go
package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/certimate-go/cli/internal/config"
)

type Client struct {
    httpClient *http.Client
    server     string
    token      string
}

func NewClient(cfg *config.Config) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        server: cfg.Server,
        token:  cfg.Token,
    }
}

func (c *Client) Do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
    var reqBody io.Reader
    if body != nil {
        jsonBytes, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("marshal request: %w", err)
        }
        reqBody = bytes.NewReader(jsonBytes)
    }

    url := c.server + path
    req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }

    req.Header.Set("Authorization", c.token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return &NetworkError{Err: err}
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        var apiErr APIError
        if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
            return fmt.Errorf("HTTP %d: failed to decode error", resp.StatusCode)
        }
        apiErr.StatusCode = resp.StatusCode
        return &apiErr
    }

    if result != nil {
        if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
            return fmt.Errorf("decode response: %w", err)
        }
    }

    return nil
}
```

### internal/api/pocketbase.go

```go
package api

import (
    "context"
    "fmt"

    "github.com/certimate-go/cli/internal/models"
)

// PocketBase ListResponse is the standard response for list queries
type ListResponse[T any] struct {
    Page       int `json:"page"`
    PerPage    int `json:"perPage"`
    TotalItems int `json:"totalItems"`
    TotalPages int `json:"totalPages"`
    Items      []T `json:"items"`
}

// ListWorkflows fetches all workflows with optional filtering
func (c *Client) ListWorkflows(ctx context.Context, filter, sort string) (*ListResponse[models.Workflow], error) {
    path := "/api/collections/workflow/records"
    params := buildQueryParams(filter, sort, 1, 100)

    var result ListResponse[models.Workflow]
    if err := c.Do(ctx, "GET", path+"?"+params, nil, &result); err != nil {
        return nil, fmt.Errorf("list workflows: %w", err)
    }
    return &result, nil
}

// GetWorkflow fetches a single workflow by ID
func (c *Client) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
    path := fmt.Sprintf("/api/collections/workflow/records/%s", id)

    var result models.Workflow
    if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
        return nil, fmt.Errorf("get workflow: %w", err)
    }
    return &result, nil
}

// RunWorkflow executes a workflow via custom API
func (c *Client) RunWorkflow(ctx context.Context, workflowID string) (*models.WorkflowRun, error) {
    path := fmt.Sprintf("/api/workflows/%s/runs", workflowID)
    body := map[string]string{"trigger": "manual"}

    var result struct {
        RunId string `json:"runId"`
    }
    if err := c.Do(ctx, "POST", path, body, &result); err != nil {
        return nil, fmt.Errorf("run workflow: %w", err)
    }

    // Fetch the created run
    return c.GetWorkflowRun(ctx, result.RunId)
}

// CancelWorkflowRun cancels a running workflow
func (c *Client) CancelWorkflowRun(ctx context.Context, workflowID, runID string) error {
    path := fmt.Sprintf("/api/workflows/%s/runs/%s/cancel", workflowID, runID)
    return c.Do(ctx, "POST", path, nil, nil)
}
```

### internal/models/workflow.go

```go
package models

import "time"

type WorkflowTriggerType string

const (
    WorkflowTriggerManual    WorkflowTriggerType = "manual"
    WorkflowTriggerScheduled WorkflowTriggerType = "scheduled"
)

type Workflow struct {
    ID          string              `json:"id"`
    Name        string              `json:"name"`
    Description string              `json:"description"`
    Trigger     WorkflowTriggerType `json:"trigger"`
    Enabled     bool                `json:"enabled"`
    Graph       WorkflowGraph       `json:"graph"`
    Created     time.Time           `json:"created"`
    Updated     time.Time           `json:"updated"`
}

type WorkflowGraph struct {
    Nodes []WorkflowNode `json:"nodes"`
    Edges []WorkflowEdge `json:"edges"`
}

type WorkflowNode struct {
    ID   string                 `json:"id"`
    Type string                 `json:"type"`
    Data map[string]interface{} `json:"data"`
}

type WorkflowEdge struct {
    ID     string `json:"id"`
    Source string `json:"source"`
    Target string `json:"target"`
}
```

### cmd/workflow.go

```go
package cmd

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "github.com/certimate-go/cli/internal/api"
    "github.com/certimate-go/cli/internal/config"
    "github.com/certimate-go/cli/internal/output"
)

var workflowCmd = &cobra.Command{
    Use:   "workflow",
    Short: "Manage certificate workflows",
}

var workflowListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all workflows",
    RunE:  runWorkflowList,
}

var workflowGetCmd = &cobra.Command{
    Use:   "get WORKFLOW_ID",
    Short: "Get workflow details",
    Args:  cobra.ExactArgs(1),
    RunE:  runWorkflowGet,
}

var workflowRunCmd = &cobra.Command{
    Use:   "run WORKFLOW_ID",
    Short: "Execute a workflow",
    Args:  cobra.ExactArgs(1),
    RunE:  runWorkflowRun,
}

var waitFlag bool

func init() {
    rootCmd.AddCommand(workflowCmd)
    workflowCmd.AddCommand(workflowListCmd)
    workflowCmd.AddCommand(workflowGetCmd)
    workflowCmd.AddCommand(workflowRunCmd)

    workflowRunCmd.Flags().BoolVar(&waitFlag, "wait", false, "wait for workflow to complete")
}

func runWorkflowList(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(profile)
    if err != nil {
        return fmt.Errorf("load config: %w", err)
    }

    client := api.NewClient(cfg)
    resp, err := client.ListWorkflows(context.Background(), "", "-created")
    if err != nil {
        return err
    }

    return output.Print(resp.Items, outputFmt)
}

func runWorkflowGet(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(profile)
    if err != nil {
        return fmt.Errorf("load config: %w", err)
    }

    client := api.NewClient(cfg)
    workflow, err := client.GetWorkflow(context.Background(), args[0])
    if err != nil {
        return err
    }

    return output.Print(workflow, outputFmt)
}

func runWorkflowRun(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(profile)
    if err != nil {
        return fmt.Errorf("load config: %w", err)
    }

    client := api.NewClient(cfg)
    run, err := client.RunWorkflow(context.Background(), args[0])
    if err != nil {
        return err
    }

    result := map[string]interface{}{
        "status":  "started",
        "run_id":  run.ID,
        "message": "Workflow execution started",
    }

    if waitFlag {
        // Poll for completion
        result, err = waitForCompletion(client, run.ID)
        if err != nil {
            return err
        }
    }

    return output.Print(result, outputFmt)
}
```

### skills/ctm-workflow/SKILL.md

```yaml
---
name: ctm-workflow
version: 1.0.0
description: |
  Certimate CLI: Workflow management for SSL certificate automation.
  Use when: (1) Listing or viewing certificate workflows, (2) Executing
  certificate issuance/renewal workflows, (3) Checking workflow execution
  status or history, (4) Canceling running workflows.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---
# Certimate Workflow Management

Manage SSL certificate workflows through the command line.

## Prerequisites

1. Install certimate CLI
2. Configure server connection:
   ```bash
   certimate config set --server https://certimate.example.com --token YOUR_TOKEN
   ```

## Available Commands

### List Workflows

```bash
certimate workflow list [--output json|table]
```

**Example Output (JSON):**
```json
[
  {
    "id": "abc123",
    "name": "Production Certificate",
    "trigger": "scheduled",
    "enabled": true,
    "created": "2026-03-01T00:00:00Z"
  }
]
```

### Get Workflow Details

```bash
certimate workflow get WORKFLOW_ID
```

### Execute Workflow

```bash
# Fire-and-forget
certimate workflow run WORKFLOW_ID

# Wait for completion
certimate workflow run WORKFLOW_ID --wait
```

**Example Output:**
```json
{
  "status": "started",
  "run_id": "run_abc123",
  "message": "Workflow execution started"
}
```

### Cancel Workflow

```bash
certimate workflow cancel WORKFLOW_ID RUN_ID
```

## Error Handling

| Exit Code | Meaning | Action |
|-----------|---------|--------|
| 3 | Authentication error | Re-run `certimate config set` |
| 4 | Network error | Check server connectivity |
| 5 | Not found | Verify workflow ID |

## Common Patterns

### Renew Certificate Before Expiry

```bash
# List certificates expiring soon
certimate certificate list --filter 'expiredAt<"2026-04-01"'

# Run the corresponding workflow
certimate workflow run WORKFLOW_ID --wait
```

### Monitor Workflow Execution

```bash
# Start workflow and capture run ID
RUN_ID=$(certimate workflow run WORKFLOW_ID | jq -r '.run_id')

# Check status
certimate workflow runs WORKFLOW_ID | jq ".[] | select(.id==\"$RUN_ID\")"
```
```

## Success Metrics

- [ ] CLI can list 100+ workflows without truncation (pagination working)
- [ ] Workflow execution completes in under 30 seconds for simple cases
- [ ] JSON output parses successfully with `jq`
- [ ] Skills documentation enables AI agent to execute workflows without errors
- [ ] All commands return appropriate exit codes for scripting

## Dependencies & Risks

### Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| Go | 1.21+ | Runtime |
| spf13/cobra | latest | CLI framework |
| spf13/viper | latest | Configuration |
| Certimate Server | 0.3.x+ | Backend API |

### Risks

| Risk | Mitigation |
|------|------------|
| PocketBase API changes | Use stable API version, add integration tests |
| Token security | Document security best practices, consider keychain for v2 |
| Large dataset performance | Implement pagination, add `--limit` flag |
| Network timeouts | Add configurable timeout, retry logic |

## References & Research

### Internal References

- Certimate domain models: `/Users/liuxuanyao/work/certimate/internal/domain/`
- Certimate API routes: `/Users/liuxuanyao/work/certimate/internal/rest/routes/routes.go`
- Go critical patterns: Memory file `debugging.md`

### External References

- googleworkspace/cli: https://github.com/googleworkspace/cli
- Cobra documentation: https://github.com/spf13/cobra
- Viper documentation: https://github.com/spf13/viper
- PocketBase API: https://pocketbase.io/docs/api/

### Related Work

- Original MVP spec: `/Users/liuxuanyao/work/certimate-cli/docs/mvp/readme.md`
