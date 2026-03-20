# Workflow Templates Reference

This directory contains templates and examples for creating workflows.

## Directory Structure

```
references/
├── README.md           # This file
├── node-types.md       # Detailed node type documentation
├── templates/
│   ├── basic-apply.json
│   ├── apply-deploy.json
│   └── full-workflow.json
└── examples/
    ├── cloudflare-example.json
    └── aliyun-example.json
```

## Templates

### basic-apply.json
Simple workflow that only applies for a certificate.

**Use when**: You just need to obtain a certificate without deploying it.

### apply-deploy.json
Workflow that applies for a certificate and deploys it.

**Use when**: You need to both obtain and deploy a certificate.

### full-workflow.json
Complete workflow with error handling and notifications.

**Use when**: You need robust error handling and failure notifications.

## How to Use Templates

1. Copy the template to your working directory
2. Customize the JSON for your needs
3. Create the workflow:

```bash
certimate workflow create --name "My Workflow" --graph-draft @my-config.json
```

## Customization Checklist

When customizing a template, make sure to update:

- [ ] Workflow name and description
- [ ] Node IDs (use meaningful names)
- [ ] `bizApply` configuration:
  - `challengeType`: "dns-01" or "http-01"
  - `keyAlgorithm`: "RSA2048", "EC256", etc.
  - `skipBeforeExpiryDays`: renewal threshold
- [ ] `bizDeploy` configuration:
  - `certificateOutputNodeId`: must match the apply node ID
- [ ] `bizNotify` configuration:
  - `subject`: notification subject
  - `message`: notification body
