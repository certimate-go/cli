---
name: ctm-workflow-create
version: 1.1.0
description: |
  Certimate CLI: Workflow creation and editing with complex graphDraft configuration.
  Use when: (1) Creating new certificate workflows, (2) Editing existing workflow
  configurations, (3) Understanding workflow graphDraft structure, (4) Working with
  workflow templates and node types.

  IMPORTANT: When creating workflows, collect ALL parameters step by step BEFORE
  generating the JSON. Do not generate partial configurations.
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
---

# Workflow Creation & Editing

This skill guides you through creating and editing certificate workflows using the certimate CLI. Workflows use a complex `graphDraft` structure with nodes and edges.

## AI Agent: Step-by-Step Parameter Collection

When a user wants to create a workflow, **collect ALL required parameters step by step** before generating the JSON configuration.

> **Critical**: The field name for provider type in bizApply/bizDeploy/bizNotify is `provider` (NOT `providerType`). The provider-specific config field is `providerConfig` (NOT `extendedConfig`).

### Step 1: Basic Information

Ask the user:
- **Workflow name**: What should this workflow be called?

### Step 2: Certificate Application (bizApply)

**Required fields:**
- **Identifier** (`identifier`): `"domain"` or `"ip"` - determines whether to use domains or ipaddrs
- **Domain(s)** (`domains`): Required if identifier="domain"
  - Format: Semicolon-separated string (e.g., `"example.com"` or `"example.com;*.example.com"`)
- **IP Address(es)** (`ipaddrs`): Required if identifier="ip"
  - Format: Semicolon-separated string (e.g., `"1.2.3.4"` or `"1.2.3.4;5.6.7.8"`)
- **Contact email** (`contactEmail`): ACME contact email address
- **Challenge type** (`challengeType`): `dns-01` or `http-01`
  - DNS-01: Requires DNS provider access (recommended for wildcards)
  - HTTP-01: Requires web server access (providers: `local`, `s3`, `ssh`)
- **Provider**: Which DNS/HTTP provider for ACME challenge?
  - DNS-01 providers: `aliyun-dns`, `tencentcloud-dns`, `cloudflare`, `huaweicloud-dns`, `aws-route53`, `azure-dns`, `godaddy`, `namesilo`, `powerdns`, `rfc2136`, etc. (75+ providers)
  - HTTP-01 providers: `local`, `s3`, `ssh`
- **Access credential**: Do you have an existing access credential for the provider?
  - If yes: Ask for the access ID
  - If no: Guide user to create one first with `certimate access create`

**Optional fields (with defaults):**
- **Contact email** (`contactEmail`): Required by some CAs
- **Key source** (`keySource`): `auto` (default), `reuse` (reuse last key), `custom` (user-provided key)
- **Key algorithm** (`keyAlgorithm`): `RSA2048` (default), `RSA3072`, `RSA4096`, `RSA8192`, `EC256`, `EC384`, `EC512`
- **Key content** (`keyContent`): PEM private key (only if `keySource: "custom"`)
- **Skip before expiry** (`skipBeforeExpiryDays`): Days before expiry to skip (default: 30)
- **CA provider** (`caProvider`): `letsencrypt` (default), `letsencryptstaging`, `zerossl`, `googletrustservices`, `digicert`, `sectigo`, `sslcom`, `acmeca`, `actalisssl`, `globalsignatlas`, `litessl`
- **Validity lifetime** (`validityLifetime`): e.g. `"30d"`, `"6h"`
- **Preferred chain** (`preferredChain`): Preferred certificate chain
- **ACME profile** (`acmeProfile`): ACME Profiles Extension
- **Nameservers** (`nameservers`): Custom DNS servers
- **DNS propagation** (`dnsPropagationWait`, `dnsPropagationTimeout`, `dnsTTL`): DNS propagation tuning
- **HTTP delay** (`httpDelayWait`): HTTP challenge delay wait
- **Disable flags** (`disableCommonName`, `disableFollowCNAME`, `disableARI`): Advanced options

### Step 3: Deployment (bizDeploy) - If Needed

Ask the user:
- **Deploy required?**: Does this certificate need to be deployed somewhere?
- **If yes**, ask:
  - **Deploy provider**: Where should the certificate be deployed?
    - See `references/deploy-configs.md` for the full list (120+ providers)
  - **Access credential**: Do you have an existing access credential for this deploy target?
    - If yes: Ask for the access ID
    - If no: Guide user to create one first
  - **Provider config** (`providerConfig`): Provider-specific configuration
    - See `references/deploy-configs.md` for each provider's required fields
  - **Skip on last succeeded** (`skipOnLastSucceeded`): Skip if last deployment succeeded? (default: `true`)

### Step 4: Notifications (bizNotify) - If Needed

Ask the user:
- **Notifications required?**: Should the workflow send notifications?
- **If yes**, ask:
  - **Provider**: `email`, `dingtalkbot`, `larkbot`, `wecombot`, `slackbot`, `discordbot`, `telegrambot`, `mattermost`, `webhook`
  - **Access credential**: Do you have an existing access credential?
    - If no: Guide user to create one first
  - **Provider config** (`providerConfig`): Provider-specific configuration
    - See `references/notify-configs.md` for each provider's required fields
    - For `email`: `receiverAddress` (required), `format` (optional)
    - For bot providers (`slackbot`, `discordbot`, `telegrambot`): `channelId` or `chatId` (if not in access)
    - For `webhook`: `webhookData`, `headers`, `timeout`
  - **Subject** (`subject`): Notification subject (supports `{{$variable}}` templates)
  - **Message** (`message`): Notification body (supports `{{$variable}}` templates)
  - **Skip on all prev skipped** (`skipOnAllPrevSkipped`): Skip if all previous nodes were skipped?

### Step 5: Error Handling (tryCatch)

Ask the user:
- **Error handling**: Should the workflow include try/catch error handling?
- If yes, configure catch block with bizNotify

### Step 6: Generate and Create

After collecting ALL parameters, generate the complete workflow JSON and create it:

```bash
certimate workflow create --name "Workflow Name" --graph-draft @workflow.json
```

## Deploy Configuration Reference

For detailed deploy provider configurations, see `references/deploy-configs.md`.

**Common Deploy Providers:**

| Provider | Access Type | Key ProviderConfig Fields |
|----------|-------------|--------------------------|
| `aliyun-cdn` | aliyun | `region`, `domain`, `domainMatchPattern` |
| `tencentcloud-cdn` | tencentcloud | `domain`, `domainMatchPattern` |
| `huaweicloud-cdn` | huaweicloud | `region`, `domain` |
| `ssh` | ssh | `format`, `keyPath`, `certPath`, `postCommand` |
| `k8s-secret` | k8s | `namespace`, `secretName` |
| `local` | local | `outputPath`, `format` |

## Quick Start

### Creating a Workflow

```bash
# Create from inline JSON
certimate workflow create --name "My Workflow" --graph-draft '{"nodes":[...],"edges":[]}'

# Create from file
certimate workflow create --name "My Workflow" --graph-draft @workflow-config.json

# Create from template (copy template first)
cp ~/.config/certimate-cli/templates/basic-apply.json my-workflow.json
# Edit my-workflow.json as needed
certimate workflow create --name "My Workflow" --graph-draft @my-workflow.json
```

### Editing a Workflow

```bash
# Update name only
certimate workflow edit WORKFLOW_ID --name "New Name"

# Update graphDraft
certimate workflow edit WORKFLOW_ID --graph-draft @updated-config.json

# Update multiple fields
certimate workflow edit WORKFLOW_ID --name "New Name" --description "Updated description"
```

## graphDraft Structure

The `graphDraft` is a JSON object with a `nodes` array. **There is NO `edges` array** - connections are implicit through the nested `blocks` structure.

```json
{
  "nodes": [
    { "id": "...", "type": "start", ... },
    { "id": "...", "type": "tryCatch", "blocks": [...] },
    { "id": "...", "type": "end", ... }
  ]
}
```

### Node Structure

Each node has:
- `id` - Unique identifier (string)
- `type` - Node type (see Node Types below)
- `data` - Node configuration (object with `name` and optional `config`)
- `blocks` - Nested nodes (for container types: tryCatch, tryBlock, catchBlock)

### Sequential Execution

Nodes in the same `blocks` array execute sequentially in order. No explicit edges are needed.

### Container Nodes

- `tryCatch` - Contains `blocks` array with exactly one `tryBlock` and one `catchBlock`
- `tryBlock` - Contains `blocks` array with main execution nodes
- `catchBlock` - Contains `blocks` array with error handling nodes

## Node Types

| Type | Description | Key Config Fields |
|------|-------------|-------------------|
| `start` | Entry point | `trigger`: "manual" or "scheduled" |
| `end` | Termination | None |
| `tryCatch` | Error wrapper | Contains `tryBlock` and `catchBlock` |
| `tryBlock` | Success path | Contains nodes to execute |
| `catchBlock` | Error path | Contains error handling nodes |
| `delay` | Wait | `wait`: seconds to delay |
| `condition` | Conditional | No specific config |
| `branchBlock` | Branch | `expression`: condition expression |
| `bizApply` | Apply certificate | `domains`, `challengeType`, `provider`, `providerAccessId`, `keySource`, `keyAlgorithm` |
| `bizUpload` | Upload certificate | `source`, `certificate`, `privateKey` |
| `bizMonitor` | Monitor cert expiry | `host`, `port`, `domain` |
| `bizDeploy` | Deploy certificate | `provider`, `providerAccessId`, `certificateOutputNodeId`, `providerConfig` |
| `bizNotify` | Send notification | `provider`, `providerAccessId`, `subject`, `message` |

For detailed node configurations, see `references/node-types.md`.

## Templates

Use templates for common workflow patterns:

1. **basic-apply.json** - Simple certificate application
2. **apply-deploy.json** - Certificate application with deployment
3. **full-workflow.json** - Complete workflow with try/catch and notifications

## Common Patterns

### Pattern 1: Simple Certificate Application

```
start → bizApply → end
```

Use `templates/basic-apply.json`

### Pattern 2: Certificate with Deployment

```
start → bizApply → bizDeploy → end
```

Use `templates/apply-deploy.json`

### Pattern 3: Full Workflow with Error Handling

```
start → tryCatch {
  tryBlock: bizApply → bizDeploy
  catchBlock: bizNotify → end
} → end
```

Use `templates/full-workflow.json`

## Step-by-Step Creation

1. **Choose a template** based on your needs
2. **Copy and customize** the JSON file
3. **Update node configurations**:
   - Set certificate parameters in `bizApply`
   - Configure deployment in `bizDeploy`
   - Add notification settings in `bizNotify`
4. **Create the workflow**:
   ```bash
   certimate workflow create --name "My Workflow" --graph-draft @my-config.json
   ```
5. **Enable the workflow**:
   ```bash
   certimate workflow enable WORKFLOW_ID
   ```

## Examples

See `references/examples/` for provider-specific examples:
- `cloudflare-example.json` - Cloudflare DNS configuration
- `aliyun-example.json` - Aliyun DNS configuration

## Related Commands

- `certimate workflow list` - List all workflows
- `certimate workflow get WORKFLOW_ID` - View workflow details
- `certimate workflow enable WORKFLOW_ID` - Enable a workflow
- `certimate workflow disable WORKFLOW_ID` - Disable a workflow
- `certimate workflow delete WORKFLOW_ID` - Delete a workflow
- `certimate workflow run WORKFLOW_ID` - Execute a workflow
