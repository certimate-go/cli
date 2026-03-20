---
title: feat: Implement Workflow Management CRUD Operations
type: feat
status: completed
date: 2026-03-19
---

# feat: Implement Workflow Management CRUD Operations

## Overview

Implement complete CRUD operations (Create, Edit, Delete, Enable, Disable) for workflow management in certimate-cli. The MVP version only supports read operations (list, get, run, cancel, runs). This plan adds full workflow lifecycle management following the established access management patterns.

**Reference Documentation**: `docs/workflows/readme.md`

## Problem Statement / Motivation

The certimate-cli MVP supports these workflow commands:
- `workflow list` - List all workflows
- `workflow get WORKFLOW_ID` - Get workflow details
- `workflow run WORKFLOW_ID` - Execute workflow
- `workflow cancel WORKFLOW_ID RUN_ID` - Cancel running workflow
- `workflow runs WORKFLOW_ID` - List execution history

**Missing functionality**:
- Create new workflows
- Edit existing workflows
- Delete workflows
- Enable workflows
- Disable workflows

These operations are essential for complete workflow lifecycle management. The workflow creation is particularly complex due to the nested `graphDraft` structure with multiple node types.

## Proposed Solution

Follow the existing access management CRUD pattern to implement:

1. **API Client Methods** in `internal/api/pocketbase.go`
2. **CLI Commands** in `cmd/workflow.go`
3. **Workflow Templates** in `skills/ctm-workflow/references/`
4. **Skills Documentation**:
   - **`ctm-workflow`** - Existing skill for workflow execution (run, cancel, list, get)
   - **`ctm-workflow-create`** (NEW) - Dedicated skill for creating and editing workflows (shared complex parameters)

> **Design Decision**: Create and edit operations share the same complex `graphDraft` parameter structure. A dedicated `ctm-workflow-create` skill consolidates guidance for both operations, making it easier for AI agents to understand the workflow configuration process.

### Architecture

```mermaid
graph TD
    A[CLI Commands] --> B[API Client]
    B --> C[PocketBase API]
    C --> D[(Workflow Collection)]

    subgraph "cmd/workflow.go"
        A1[workflow create]
        A2[workflow edit]
        A3[workflow delete]
        A4[workflow enable]
        A5[workflow disable]
    end

    subgraph "internal/api/pocketbase.go"
        B1[CreateWorkflow]
        B2[UpdateWorkflow]
        B3[DeleteWorkflow]
    end

    subgraph "PocketBase Endpoints"
        C1[POST /api/collections/workflow/records]
        C2[PATCH /api/collections/workflow/records/{id}]
        C3[DELETE /api/collections/workflow/records/{id}]
    end
```

## Technical Approach

### Phase 1: API Client Implementation

**File**: `internal/api/pocketbase.go`

Add three new methods following the existing pattern:

```go
// CreateWorkflow - POST /api/collections/workflow/records
func (c *Client) CreateWorkflow(ctx context.Context, workflow *models.Workflow) (*models.Workflow, error) {
    path := "/api/collections/workflow/records"
    var result models.Workflow
    if err := c.Do(ctx, "POST", path, workflow, &result); err != nil {
        return nil, fmt.Errorf("create workflow: %w", err)
    }
    return &result, nil
}

// UpdateWorkflow - PATCH /api/collections/workflow/records/{id}
func (c *Client) UpdateWorkflow(ctx context.Context, id string, workflow *models.Workflow) (*models.Workflow, error) {
    path := fmt.Sprintf("/api/collections/workflow/records/%s", id)
    var result models.Workflow
    if err := c.Do(ctx, "PATCH", path, workflow, &result); err != nil {
        return nil, fmt.Errorf("update workflow: %w", err)
    }
    return &result, nil
}

// DeleteWorkflow - DELETE /api/collections/workflow/records/{id}
func (c *Client) DeleteWorkflow(ctx context.Context, id string) error {
    path := fmt.Sprintf("/api/collections/workflow/records/%s", id)
    return c.Do(ctx, "DELETE", path, nil, nil)
}
```

### Phase 2: CLI Commands Implementation

**File**: `cmd/workflow.go`

#### 2.1 Create Command

```go
var (
    workflowCreateName        string
    workflowCreateDescription string
    workflowCreateGraphDraft  string
    workflowCreateHasDraft    bool
)

var workflowCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new workflow",
    RunE:  runWorkflowCreate,
}

func init() {
    workflowCreateCmd.Flags().StringVar(&workflowCreateName, "name", "", "Workflow name")
    workflowCreateCmd.Flags().StringVar(&workflowCreateDescription, "description", "", "Workflow description")
    workflowCreateCmd.Flags().StringVar(&workflowCreateGraphDraft, "graph-draft", "", "Workflow graph structure (JSON string or @path/to/file.json)")
    workflowCreateCmd.Flags().BoolVar(&workflowCreateHasDraft, "has-draft", true, "Mark workflow as having draft")
    workflowCreateCmd.MarkFlagRequired("name")
    workflowCreateCmd.MarkFlagRequired("graph-draft")
}

func runWorkflowCreate(cmd *cobra.Command, args []string) error {
    client, err := getAPIClient()
    if err != nil {
        return err
    }

    // Parse graphDraft JSON
    graphDraft, err := parseGraphDraftInput(workflowCreateGraphDraft)
    if err != nil {
        return err
    }

    workflow := &models.Workflow{
        Name:        workflowCreateName,
        Description: workflowCreateDescription,
        GraphDraft:  graphDraft,
        HasDraft:    workflowCreateHasDraft,
        Enabled:     false, // New workflows start disabled
    }

    created, err := client.CreateWorkflow(context.Background(), workflow)
    if err != nil {
        return err
    }

    return output.Print(created, outputFmt)
}
```

#### 2.2 Edit Command

```go
var (
    workflowEditName        string
    workflowEditDescription string
    workflowEditGraphDraft  string
)

var workflowEditCmd = &cobra.Command{
    Use:   "edit WORKFLOW_ID",
    Short: "Edit an existing workflow",
    Args:  cobra.ExactArgs(1),
    RunE:  runWorkflowEdit,
}

func init() {
    workflowEditCmd.Flags().StringVar(&workflowEditName, "name", "", "New workflow name")
    workflowEditCmd.Flags().StringVar(&workflowEditDescription, "description", "", "New workflow description")
    workflowEditCmd.Flags().StringVar(&workflowEditGraphDraft, "graph-draft", "", "New workflow graph structure (JSON string or @path/to/file.json)")
}

func runWorkflowEdit(cmd *cobra.Command, args []string) error {
    client, err := getAPIClient()
    if err != nil {
        return err
    }

    workflowID := args[0]

    // Fetch existing workflow
    existing, err := client.GetWorkflow(context.Background(), workflowID)
    if err != nil {
        return err
    }

    // Apply updates (only if flag provided)
    if workflowEditName != "" {
        existing.Name = workflowEditName
    }
    if cmd.Flags().Changed("description") {
        existing.Description = workflowEditDescription
    }
    if workflowEditGraphDraft != "" {
        graphDraft, err := parseGraphDraftInput(workflowEditGraphDraft)
        if err != nil {
            return err
        }
        existing.GraphDraft = graphDraft
    }

    updated, err := client.UpdateWorkflow(context.Background(), workflowID, existing)
    if err != nil {
        return err
    }

    return output.Print(updated, outputFmt)
}
```

#### 2.3 Delete Command

```go
var workflowDeleteForce bool

var workflowDeleteCmd = &cobra.Command{
    Use:   "delete WORKFLOW_ID",
    Short: "Delete a workflow",
    Args:  cobra.ExactArgs(1),
    RunE:  runWorkflowDelete,
}

func init() {
    workflowDeleteCmd.Flags().BoolVar(&workflowDeleteForce, "force", false, "Skip confirmation prompt")
}

func runWorkflowDelete(cmd *cobra.Command, args []string) error {
    client, err := getAPIClient()
    if err != nil {
        return err
    }

    workflowID := args[0]

    // Optional: Fetch workflow to show what will be deleted
    if !workflowDeleteForce {
        existing, err := client.GetWorkflow(context.Background(), workflowID)
        if err != nil {
            return err
        }
        fmt.Printf("Will delete workflow: %s (%s)\n", existing.Name, workflowID)
        fmt.Print("Confirm? (y/N): ")
        var response string
        fmt.Scanln(&response)
        if strings.ToLower(response) != "y" {
            fmt.Println("Cancelled")
            return nil
        }
    }

    if err := client.DeleteWorkflow(context.Background(), workflowID); err != nil {
        return err
    }

    fmt.Printf("Workflow %s deleted\n", workflowID)
    return nil
}
```

#### 2.4 Enable Command

```go
var workflowEnableCmd = &cobra.Command{
    Use:   "enable WORKFLOW_ID",
    Short: "Enable a workflow",
    Args:  cobra.ExactArgs(1),
    RunE:  runWorkflowEnable,
}

func runWorkflowEnable(cmd *cobra.Command, args []string) error {
    client, err := getAPIClient()
    if err != nil {
        return err
    }

    workflowID := args[0]

    // Fetch existing workflow
    existing, err := client.GetWorkflow(context.Background(), workflowID)
    if err != nil {
        return err
    }

    if existing.Enabled {
        fmt.Printf("Workflow %s is already enabled\n", workflowID)
        return nil
    }

    existing.Enabled = true

    updated, err := client.UpdateWorkflow(context.Background(), workflowID, existing)
    if err != nil {
        return err
    }

    return output.Print(updated, outputFmt)
}
```

#### 2.5 Disable Command

```go
var workflowDisableCmd = &cobra.Command{
    Use:   "disable WORKFLOW_ID",
    Short: "Disable a workflow",
    Args:  cobra.ExactArgs(1),
    RunE:  runWorkflowDisable,
}

func runWorkflowDisable(cmd *cobra.Command, args []string) error {
    client, err := getAPIClient()
    if err != nil {
        return err
    }

    workflowID := args[0]

    // Fetch existing workflow
    existing, err := client.GetWorkflow(context.Background(), workflowID)
    if err != nil {
        return err
    }

    if !existing.Enabled {
        fmt.Printf("Workflow %s is already disabled\n", workflowID)
        return nil
    }

    existing.Enabled = false

    updated, err := client.UpdateWorkflow(context.Background(), workflowID, existing)
    if err != nil {
        return err
    }

    return output.Print(updated, outputFmt)
}
```

#### 2.6 Helper Functions

```go
// parseGraphDraftInput parses graphDraft from JSON string or file path
func parseGraphDraftInput(input string) (models.WorkflowGraph, error) {
    var data []byte
    var err error

    if filePath, ok := strings.CutPrefix(input, "@"); ok {
        data, err = os.ReadFile(filePath)
        if err != nil {
            return models.WorkflowGraph{}, fmt.Errorf("failed to read graph-draft file: %w", err)
        }
    } else {
        data = []byte(input)
    }

    var graphDraft models.WorkflowGraph
    if err := json.Unmarshal(data, &graphDraft); err != nil {
        return models.WorkflowGraph{}, fmt.Errorf("invalid graph-draft JSON format: %w", err)
    }

    return graphDraft, nil
}
```

### Phase 3: Workflow Templates

Create template files for common workflow patterns in `skills/ctm-workflow/references/`:

#### 3.1 Directory Structure

```
skills/
├── ctm-workflow/
│   └── SKILL.md                     # Execution skill (existing, update with CRUD refs)
└── ctm-workflow-create/             # NEW: Dedicated create/edit skill
    ├── SKILL.md                     # Main skill documentation
    └── references/
        ├── README.md                # Template overview
        ├── node-types.md            # Node type documentation
        ├── templates/
        │   ├── basic-apply.json     # Basic certificate application
        │   ├── apply-deploy.json    # Certificate with deployment
        │   └── full-workflow.json   # Full workflow with try/catch
        └── examples/
            ├── cloudflare-example.json
            └── aliyun-example.json
```

#### 3.2 Node Types Reference (`references/node-types.md`)

| Node Type | Description | Required Config Fields |
|-----------|-------------|----------------------|
| `start` | Workflow entry point | `trigger` (manual/scheduled) |
| `tryCatch` | Error handling wrapper | - |
| `tryBlock` | Main execution block | - |
| `catchBlock` | Error handling block | - |
| `bizApply` | Certificate application | `challengeType`, `keySource`, `keyAlgorithm`, `skipBeforeExpiryDays` |
| `bizDeploy` | Certificate deployment | `skipOnLastSucceeded`, `certificateOutputNodeId` |
| `bizNotify` | Notification | `subject`, `message` |
| `end` | Workflow termination | - |

#### 3.3 Basic Template (`templates/basic-apply.json`)

```json
{
  "nodes": [
    {
      "id": "start-node",
      "type": "start",
      "data": {
        "name": "Start",
        "config": {
          "trigger": "manual"
        }
      }
    },
    {
      "id": "apply-node",
      "type": "bizApply",
      "data": {
        "name": "Apply Certificate",
        "config": {
          "challengeType": "dns-01",
          "keySource": "auto",
          "keyAlgorithm": "RSA2048",
          "skipBeforeExpiryDays": 30
        }
      }
    },
    {
      "id": "end-node",
      "type": "end",
      "data": {
        "name": "End"
      }
    }
  ],
  "edges": [
    {"source": "start-node", "target": "apply-node"},
    {"source": "apply-node", "target": "end-node"}
  ]
}
```

#### 3.4 Full Workflow Template (`templates/full-workflow.json`)

```json
{
  "nodes": [
    {
      "id": "start-node",
      "type": "start",
      "data": {
        "name": "Start",
        "config": {
          "trigger": "manual"
        }
      }
    },
    {
      "id": "try-catch-node",
      "type": "tryCatch",
      "data": {
        "name": "Try Execute"
      },
      "blocks": [
        {
          "id": "try-block",
          "type": "tryBlock",
          "data": {"name": ""},
          "blocks": [
            {
              "id": "apply-node",
              "type": "bizApply",
              "data": {
                "name": "Apply Certificate",
                "config": {
                  "challengeType": "dns-01",
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
                "name": "Deploy Certificate",
                "config": {
                  "skipOnLastSucceeded": true,
                  "certificateOutputNodeId": "apply-node"
                }
              }
            }
          ]
        },
        {
          "id": "catch-block",
          "type": "catchBlock",
          "data": {"name": "On Error"},
          "blocks": [
            {
              "id": "notify-node",
              "type": "bizNotify",
              "data": {
                "name": "Send Notification",
                "config": {
                  "subject": "[Certimate] Workflow Failure Alert!",
                  "message": "Your workflow \"{{ $workflow.name }}\" run has failed."
                }
              }
            },
            {
              "id": "end-error-node",
              "type": "end",
              "data": {"name": "End"}
            }
          ]
        }
      ]
    },
    {
      "id": "end-success-node",
      "type": "end",
      "data": {"name": "End"}
    }
  ]
}
```

### Phase 4: Skills Documentation

#### 4.1 New Skill: `ctm-workflow-create`

Create a dedicated skill for workflow creation and editing since both operations share the complex `graphDraft` parameter.

**File**: `skills/ctm-workflow-create/SKILL.md`

```yaml
---
name: ctm-workflow-create
version: 1.0.0
description: |
  Certimate CLI: Workflow creation and editing with complex graphDraft configuration.
  Use when: (1) Creating new certificate workflows, (2) Editing existing workflow
  configurations, (3) Understanding workflow graphDraft structure, (4) Working with
  workflow templates and node types.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---
```

**Key sections in the skill**:

1. **Workflow Configuration Overview**: Explain graphDraft structure
2. **Node Types Reference**: Document all available node types
3. **Creating Workflows**: Step-by-step guide with examples
4. **Editing Workflows**: How to update existing configurations
5. **Templates**: Available templates and when to use them
6. **Common Patterns**: Best practices for workflow design

#### 4.2 Update Existing Skill: `ctm-workflow`

**File**: `skills/ctm-workflow/SKILL.md`

Update the existing skill to focus on execution and add CRUD operation references:

```markdown
## Workflow CRUD Operations

For creating and editing workflows with complex graphDraft configuration, use the
dedicated `ctm-workflow-create` skill.

### Quick Reference

- **Create/Edit Workflows**: See `ctm-workflow-create` skill
- **Delete Workflow**: `certimate workflow delete WORKFLOW_ID [--force]`
- **Enable Workflow**: `certimate workflow enable WORKFLOW_ID`
- **Disable Workflow**: `certimate workflow disable WORKFLOW_ID`
```

## System-Wide Impact

### Interaction Graph

- **Create**: CLI → API Client → PocketBase POST → Workflow Collection
- **Edit**: CLI → API Client → PocketBase GET → Merge → PocketBase PATCH → Workflow Collection
- **Delete**: CLI → API Client → PocketBase DELETE → Workflow Collection (cascades to runs?)
- **Enable/Disable**: CLI → API Client → PocketBase GET → Update enabled → PocketBase PATCH

### Error Propagation

| Error Type | Exit Code | User Message |
|------------|-----------|--------------|
| Invalid JSON | 2 | "invalid graph-draft JSON format: ..." |
| File not found | 2 | "failed to read graph-draft file: ..." |
| Auth failure | 3 | "authentication failed: ..." |
| Network error | 4 | "network error: ..." |
| Not found | 5 | "workflow not found: ..." |
| Validation error | 1 | Server validation message |

### State Lifecycle Risks

- **Partial Updates**: Edit operations use fetch-merge-update pattern; if merge fails, no update occurs
- **Concurrent Edits**: Last write wins; no conflict detection
- **Active Runs**: Deleting/editing a workflow with active runs may cause undefined behavior

## Acceptance Criteria

### Functional Requirements

- [x] **Create Workflow**: `certimate workflow create --name <name> --graph-draft <json>`
  - Accepts JSON string or file path (`@path/to/file.json`)
  - Returns created workflow with ID
  - Validates required flags (name, graph-draft)
  - New workflows start disabled (`enabled: false`)

- [x] **Edit Workflow**: `certimate workflow edit WORKFLOW_ID [--name <name>] [--description <desc>] [--graph-draft <json>]`
  - Fetches existing workflow and merges updates
  - Supports partial updates (only specified fields are changed)
  - Returns updated workflow

- [x] **Delete Workflow**: `certimate workflow delete WORKFLOW_ID [--force]`
  - Prompts for confirmation unless `--force` flag
  - Deletes workflow by ID
  - Returns confirmation message

- [x] **Enable Workflow**: `certimate workflow enable WORKFLOW_ID`
  - Sets `enabled: true`
  - Returns updated workflow
  - Handles already-enabled case gracefully

- [x] **Disable Workflow**: `certimate workflow disable WORKFLOW_ID`
  - Sets `enabled: false`
  - Returns updated workflow
  - Handles already-disabled case gracefully

### Non-Functional Requirements

- [x] **Error Handling**: Proper error messages with exit codes (0-5)
- [x] **Output Formats**: Support JSON and table output formats
- [x] **Documentation**: Skills documentation updated with CRUD examples

### Quality Gates

- [x] CLI compiles without errors
- [x] Commands display proper help text
- [x] API methods follow existing patterns
- [x] Templates directory created with sample templates

## Success Metrics

- All CRUD commands work correctly against live Certimate instance
- Commands follow consistent CLI patterns with access management
- Templates provide useful starting points for common workflows
- Documentation enables AI agents to create workflows effectively

## Dependencies & Prerequisites

- **Certimate Server**: Running instance with PocketBase backend
- **Authentication**: Valid API token configured
- **Go Version**: 1.24.13 or compatible

## Risk Analysis & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Complex graphDraft structure confuses users | High | Provide templates and examples |
| Server validation differs from CLI assumptions | Medium | Rely on server validation; show clear errors |
| Concurrent edit conflicts | Low | Document last-write-wins behavior |
| Deleting workflow with active runs | Medium | Add warning in documentation |

## Implementation Plan

### Files to Modify

1. **`internal/api/pocketbase.go`**
   - Add `CreateWorkflow`, `UpdateWorkflow`, `DeleteWorkflow` methods
   - Lines: ~50 new lines

2. **`cmd/workflow.go`**
   - Add `workflowCreateCmd`, `workflowEditCmd`, `workflowDeleteCmd`, `workflowEnableCmd`, `workflowDisableCmd`
   - Add `parseGraphDraftInput` helper
   - Register commands in `init()`
   - Lines: ~200 new lines

3. **`skills/ctm-workflow/SKILL.md`**
   - Add CRUD documentation section
   - Add template usage examples
   - Lines: ~50 new lines

### Files to Create

1. **`skills/ctm-workflow-create/SKILL.md`** - New dedicated skill for create/edit operations
2. **`skills/ctm-workflow-create/references/README.md`** - Template overview
3. **`skills/ctm-workflow-create/references/node-types.md`** - Node type documentation
4. **`skills/ctm-workflow-create/references/templates/basic-apply.json`** - Basic template
5. **`skills/ctm-workflow-create/references/templates/apply-deploy.json`** - Apply + Deploy template
6. **`skills/ctm-workflow-create/references/templates/full-workflow.json`** - Full workflow template

## Future Considerations

### Potential Enhancements

1. **Workflow Validation Command**: `certimate workflow validate --graph-draft @config.json`
2. **Template System**: `certimate workflow create --template basic-apply`
3. **Workflow Clone**: `certimate workflow clone SOURCE_ID --name "Copy"`
4. **Export/Import**: `certimate workflow export ID > backup.json`
5. **--dry-run Flag**: Validate without creating

### Extensibility

- New node types can be added without CLI changes
- Templates are JSON files, easily extensible
- Command structure allows additional flags

## Sources & References

### Internal References

- **Access Management Pattern**: `docs/plans/2026-03-19-feat-access-management-implementation-plan.md`
- **Workflow Model**: `internal/models/workflow.go`
- **API Client Pattern**: `internal/api/pocketbase.go:171-229`
- **Existing Workflow Commands**: `cmd/workflow.go`
- **Skills Pattern**: `skills/ctm-access/SKILL.md`

### External References

- **PocketBase API**: https://pocketbase.io/docs/api-records/
- **Cobra CLI Framework**: https://github.com/spf13/cobra
- **Google Workspace CLI** (reference organization): https://github.com/googleworkspace/cli
- **Certimate Documentation**: https://deepwiki.com/certimate-go/certimate

### Related Documentation

- **Workflow Docs**: `docs/workflows/readme.md`
- **MVP Plan**: `docs/plans/2026-03-17-feat-certimate-cli-mvp-plan.md`

## Appendix: Example Commands

### Create Workflow Examples

```bash
# Create with inline JSON (simple)
certimate workflow create --name "My Workflow" --graph-draft '{"nodes":[{"id":"start","type":"start","data":{"name":"Start","config":{"trigger":"manual"}}},{"id":"end","type":"end","data":{"name":"End"}}],"edges":[]}'

# Create from file
certimate workflow create --name "Production SSL" --graph-draft @./workflows/prod-ssl.json

# Create with description
certimate workflow create --name "My Workflow" --description "Auto-renew SSL for example.com" --graph-draft @template.json
```

### Edit Workflow Examples

```bash
# Rename workflow
certimate workflow edit abc123 --name "New Name"

# Update description
certimate workflow edit abc123 --description "Updated description"

# Update graphDraft
certimate workflow edit abc123 --graph-draft @updated-workflow.json
```

### Enable/Disable Examples

```bash
# Enable workflow
certimate workflow enable abc123

# Disable workflow
certimate workflow disable abc123
```

### Delete Workflow Examples

```bash
# Delete with confirmation
certimate workflow delete abc123

# Force delete
certimate workflow delete abc123 --force
```
