---
name: ctm-workflow
version: 1.1.0
description: |
  Certimate CLI: Workflow management for SSL certificate automation.
  Use when: (1) Listing or viewing certificate workflows, (2) Executing
  certificate issuance/renewal workflows, (3) Checking workflow execution
  status or history, (4) Canceling running workflows, (5) Enabling/disabling
  or deleting workflows.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---
# Certimate Workflow Management

Manage SSL certificate workflows through the command line.

## Available Commands

### List Workflows

```bash
certimate workflow list [--output json|table] [--limit N]
```

**Example Output (JSON):**
```json
{
  "page": 1,
  "perPage": 100,
  "totalItems": 2,
  "items": [
    {
      "id": "abc123",
      "name": "Production Certificate",
      "trigger": "scheduled",
      "enabled": true,
      "created": "2026-03-01T00:00:00Z"
    }
  ]
}
```

### Get Workflow Details

```bash
certimate workflow get WORKFLOW_ID
```

### Execute Workflow

```bash
# Fire-and-forget (returns immediately with run ID)
certimate workflow run WORKFLOW_ID

# Wait for completion
certimate workflow run WORKFLOW_ID --wait

# With timeout (default 300 seconds)
certimate workflow run WORKFLOW_ID --wait --timeout 600
```

**Example Output:**
```json
{
  "status": "started",
  "run_id": "run_abc123",
  "workflow_id": "wf_xyz",
  "message": "Workflow execution started"
}
```

### Cancel Workflow

```bash
certimate workflow cancel WORKFLOW_ID RUN_ID
```

### List Workflow Runs

```bash
certimate workflow runs WORKFLOW_ID [--limit N]
```

## Common Patterns

### Renew Certificate Before Expiry

```bash
# 1. Find certificate expiring soon
certimate certificate list --output table

# 2. Find the workflow for that certificate
certimate workflow list

# 3. Run the workflow
certimate workflow run WORKFLOW_ID --wait
```

### Monitor Workflow Execution

```bash
# Start workflow and capture run ID
certimate workflow run WORKFLOW_ID | jq -r '.run_id'

# Check run history
certimate workflow runs WORKFLOW_ID
```

### Bulk Workflow Execution

```bash
# Run all enabled workflows
for id in $(certimate workflow list | jq -r '.items[] | select(.enabled == true) | .id'); do
  certimate workflow run "$id"
done
```

## Error Handling

| Exit Code | Meaning | Action |
|-----------|---------|--------|
| 3 | Authentication error | Re-run `certimate config set` |
| 4 | Network error | Check server connectivity |
| 5 | Not found | Verify workflow ID |

## Workflow CRUD Operations

### Create/Edit Workflows

For creating and editing workflows with complex `graphDraft` configuration, use the dedicated `ctm-workflow-create` skill which provides:
- Detailed graphDraft structure documentation
- Node type references
- Ready-to-use templates
- Step-by-step creation guides

### Quick CRUD Reference

```bash
# Create workflow (see ctm-workflow-create skill for details)
certimate workflow create --name "My Workflow" --graph-draft @config.json

# Edit workflow
certimate workflow edit WORKFLOW_ID [--name "New Name"] [--description "Desc"] [--graph-draft @config.json]

# Delete workflow (prompts for confirmation)
certimate workflow delete WORKFLOW_ID

# Delete without confirmation
certimate workflow delete WORKFLOW_ID --force

# Enable workflow
certimate workflow enable WORKFLOW_ID

# Disable workflow
certimate workflow disable WORKFLOW_ID
```

### Related Skills

- **ctm-workflow-create** - Detailed workflow creation and editing with templates
- **ctm-access** - Managing provider access credentials
- **ctm-certificate** - Viewing and downloading certificates
