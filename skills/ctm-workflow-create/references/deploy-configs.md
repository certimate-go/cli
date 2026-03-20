# Deploy Configuration Reference

This document provides detailed configuration for all deploy providers. Each provider requires specific `providerConfig` fields in addition to the access credential.

> **Important**: The field names are `provider` (not `providerType`) and `providerConfig` (not `extendedConfig`). This matches the Go struct tags in the certimate source code.

## How Deploy Configuration Works

In the workflow `graphDraft`, a `bizDeploy` node has this structure:

```json
{
  "id": "deploy-node",
  "type": "bizDeploy",
  "data": {
    "name": "Deploy Certificate",
    "config": {
      "provider": "aliyun-cdn",
      "providerAccessId": "ACCESS_ID",
      "skipOnLastSucceeded": true,
      "certificateOutputNodeId": "apply-node",
      "providerConfig": {
        "region": "cn-hangzhou",
        "domain": "cdn.example.com",
        "domainMatchPattern": "exact"
      }
    }
  }
}
```

**Key Fields:**
- `provider`: The deploy provider type (determines which providerConfig is needed)
- `providerAccessId`: ID of the access credential (created via `certimate access create`)
- `skipOnLastSucceeded`: Skip deployment if last run succeeded (default: `true`)
- `certificateOutputNodeId`: ID of the `bizApply` or `bizUpload` node that outputs the certificate
- `providerConfig`: Provider-specific configuration (see below)

---

## Aliyun (aliyun)

### aliyun-alb (Application Load Balancer)
```json
{
  "region": "cn-hangzhou",
  "loadbalancerId": "alb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### aliyun-apigw (API Gateway)
```json
{
  "region": "cn-hangzhou",
  "groupId": "xxx",
  "domain": "api.example.com"
}
```

### aliyun-cas (Certificate Management Service)
```json
{
  "region": "cn-hangzhou"
}
```

### aliyun-casdeploy (CAS Deploy)
```json
{
  "region": "cn-hangzhou",
  "resourceGroupId": "rg-xxx",
  "contactIds": ["contact-1", "contact-2"]
}
```

### aliyun-cdn (Content Delivery Network)
```json
{
  "region": "cn-hangzhou",
  "domain": "cdn.example.com",
  "domainMatchPattern": "exact"
}
```
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| region | string | No | Aliyun region (e.g., cn-hangzhou, ap-southeast-1) |
| domain | string | Yes | Domain to deploy certificate to |
| domainMatchPattern | string | No | Pattern: `exact`, `wildcard`, `certsan` |

### aliyun-clb (Classic Load Balancer)
```json
{
  "region": "cn-hangzhou",
  "loadbalancerId": "lb-xxx",
  "listenerPort": 443,
  "domain": "example.com"
}
```

### aliyun-dcdn (Dynamic Route for CDN)
```json
{
  "region": "cn-hangzhou",
  "domain": "dcdn.example.com",
  "domainMatchPattern": "exact"
}
```

### aliyun-ddospro (DDoS Protection)
```json
{
  "region": "cn-hangzhou",
  "domain": "example.com"
}
```

### aliyun-esa (Edge Security Acceleration)
```json
{
  "region": "cn-hangzhou",
  "siteId": "site-xxx",
  "domain": "example.com"
}
```

### aliyun-esasaas (ESA SaaS)
```json
{
  "region": "cn-hangzhou",
  "siteId": "site-xxx",
  "domain": "example.com"
}
```

### aliyun-fc (Function Compute)
```json
{
  "region": "cn-hangzhou",
  "domain": "example.com"
}
```

### aliyun-ga (Global Acceleration)
```json
{
  "region": "cn-hangzhou",
  "acceleratorId": "ga-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### aliyun-live (Live Streaming)
```json
{
  "region": "cn-hangzhou",
  "domain": "live.example.com",
  "domainMatchPattern": "exact"
}
```

### aliyun-nlb (Network Load Balancer)
```json
{
  "region": "cn-hangzhou",
  "loadbalancerId": "nlb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### aliyun-oss (Object Storage Service)
```json
{
  "region": "cn-hangzhou",
  "bucket": "my-bucket",
  "domain": "oss.example.com"
}
```

### aliyun-vod (Video on Demand)
```json
{
  "region": "cn-hangzhou",
  "domain": "vod.example.com"
}
```

### aliyun-waf (Web Application Firewall)
```json
{
  "region": "cn-hangzhou",
  "instanceId": "waf-xxx",
  "domain": "example.com"
}
```

---

## Tencent Cloud (tencentcloud)

### tencentcloud-cdn (Content Delivery Network)
```json
{
  "domain": "cdn.example.com",
  "domainMatchPattern": "exact"
}
```

### tencentcloud-clb (Cloud Load Balancer)
```json
{
  "region": "ap-guangzhou",
  "loadbalancerId": "lb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### tencentcloud-cos (Cloud Object Storage)
```json
{
  "region": "ap-guangzhou",
  "bucket": "my-bucket",
  "domain": "cos.example.com"
}
```

### tencentcloud-css (Cloud Streaming Services)
```json
{
  "domain": "live.example.com"
}
```

### tencentcloud-ecdn (Enterprise CDN)
```json
{
  "domain": "ecdn.example.com"
}
```

### tencentcloud-eo (EdgeOne)
```json
{
  "zoneId": "zone-xxx",
  "domain": "example.com"
}
```

### tencentcloud-gaap (Global Application Acceleration)
```json
{
  "proxyId": "proxy-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### tencentcloud-scf (Serverless Cloud Function)
```json
{
  "region": "ap-guangzhou",
  "domain": "example.com"
}
```

### tencentcloud-ssl (SSL Certificate Service)
```json
{
  "region": "ap-guangzhou"
}
```

### tencentcloud-ssldeploy (SSL Deploy)
```json
{
  "region": "ap-guangzhou",
  "cloudCertId": "cert-xxx"
}
```

### tencentcloud-sslupdate (SSL Update)
```json
{
  "region": "ap-guangzhou",
  "cloudCertId": "cert-xxx"
}
```

### tencentcloud-vod (Video on Demand)
```json
{
  "subAppId": "123456",
  "domain": "vod.example.com"
}
```

### tencentcloud-waf (Web Application Firewall)
```json
{
  "region": "ap-guangzhou",
  "instanceId": "waf-xxx",
  "domain": "example.com"
}
```

---

## AWS (aws)

### aws-acm (AWS Certificate Manager)
```json
{
  "region": "us-east-1"
}
```

### aws-cloudfront (CloudFront)
```json
{
  "region": "us-east-1",
  "distributionId": "E1234567890ABC"
}
```

### aws-iam (Identity and Access Management)
```json
{
  "region": "us-east-1"
}
```

---

## Huawei Cloud (huaweicloud)

### huaweicloud-cdn (Content Delivery Network)
```json
{
  "region": "cn-north-4",
  "domain": "cdn.example.com"
}
```

### huaweicloud-elb (Elastic Load Balance)
```json
{
  "region": "cn-north-4",
  "loadbalancerId": "lb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### huaweicloud-obs (Object Storage Service)
```json
{
  "region": "cn-north-4",
  "bucket": "my-bucket",
  "domain": "obs.example.com"
}
```

### huaweicloud-scm (SSL Certificate Manager)
```json
{
  "region": "cn-north-4"
}
```

### huaweicloud-waf (Web Application Firewall)
```json
{
  "region": "cn-north-4",
  "instanceId": "waf-xxx",
  "domain": "example.com"
}
```

---

## VolcEngine (volcengine)

### volcengine-alb (Application Load Balancer)
```json
{
  "region": "cn-beijing",
  "loadbalancerId": "alb-xxx",
  "listenerId": "lsr-xxx"
}
```

### volcengine-cdn (Content Delivery Network)
```json
{
  "domain": "cdn.example.com"
}
```

### volcengine-certcenter (Certificate Center)
```json
{
  "region": "cn-beijing"
}
```

### volcengine-clb (Cloud Load Balancer)
```json
{
  "region": "cn-beijing",
  "loadbalancerId": "clb-xxx",
  "listenerId": "lsr-xxx"
}
```

### volcengine-dcdn (Dynamic CDN)
```json
{
  "domain": "dcdn.example.com"
}
```

### volcengine-imagex (ImageX)
```json
{
  "serviceId": "service-xxx",
  "domain": "image.example.com"
}
```

### volcengine-live (Live Streaming)
```json
{
  "domain": "live.example.com"
}
```

### volcengine-tos (Tinder Object Storage)
```json
{
  "region": "cn-beijing",
  "bucket": "my-bucket",
  "domain": "tos.example.com"
}
```

### volcengine-vod (Video on Demand)
```json
{
  "spaceName": "my-space",
  "domain": "vod.example.com"
}
```

### volcengine-waf (Web Application Firewall)
```json
{
  "region": "cn-beijing",
  "instanceId": "waf-xxx",
  "domain": "example.com"
}
```

---

## JD Cloud (jdcloud)

### jdcloud-alb (Application Load Balancer)
```json
{
  "region": "cn-north-1",
  "loadbalancerId": "alb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### jdcloud-cdn (Content Delivery Network)
```json
{
  "domain": "cdn.example.com"
}
```

### jdcloud-live (Live Streaming)
```json
{
  "domain": "live.example.com"
}
```

### jdcloud-vod (Video on Demand)
```json
{
  "domain": "vod.example.com"
}
```

---

## Qiniu (qiniu)

### qiniu-cdn (Content Delivery Network)
```json
{
  "domain": "cdn.example.com"
}
```

### qiniu-kodo (Object Storage)
```json
{
  "bucket": "my-bucket",
  "domain": "kodo.example.com"
}
```

### qiniu-pili (Live Streaming)
```json
{
  "hub": "my-hub",
  "domain": "live.example.com"
}
```

---

## Baidu Cloud (baiducloud)

### baiducloud-appblb (Application BLB)
```json
{
  "region": "bj",
  "loadbalancerId": "lb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### baiducloud-blb (BLB)
```json
{
  "region": "bj",
  "loadbalancerId": "lb-xxx",
  "listenerPort": 443,
  "domain": "example.com"
}
```

### baiducloud-cdn (Content Delivery Network)
```json
{
  "domain": "cdn.example.com"
}
```

### baiducloud-cert (Certificate Service)
```json
{
  "region": "bj"
}
```

---

## UCloud (ucloud)

### ucloud-ualb (Application Load Balancer)
```json
{
  "region": "cn-bj2",
  "loadbalancerId": "alb-xxx",
  "listenerId": "lsr-xxx",
  "domain": "example.com"
}
```

### ucloud-ucdn (Content Delivery Network)
```json
{
  "domain": "cdn.example.com"
}
```

### ucloud-uclb (Cloud Load Balancer)
```json
{
  "region": "cn-bj2",
  "loadbalancerId": "lb-xxx",
  "listenerPort": 443,
  "domain": "example.com"
}
```

### ucloud-uewaf (Edge WAF)
```json
{
  "instanceId": "waf-xxx",
  "domain": "example.com"
}
```

### ucloud-upathx (PathX)
```json
{
  "instanceId": "pathx-xxx",
  "domain": "example.com"
}
```

### ucloud-us3 (Object Storage)
```json
{
  "bucket": "my-bucket",
  "domain": "us3.example.com"
}
```

---

## SSH (ssh)

### ssh (SSH/SCP Deployment)

This is a versatile deployer for uploading certificates to remote servers via SSH.

```json
{
  "useSCP": false,
  "format": "PEM",
  "keyPath": "/etc/nginx/ssl/example.com.key",
  "certPath": "/etc/nginx/ssl/example.com.crt",
  "certPathForServerOnly": "/etc/nginx/ssl/example.com.server.crt",
  "certPathForIntermediaOnly": "/etc/nginx/ssl/example.com.chain.crt",
  "pfxPassword": "",
  "jksAlias": "",
  "jksKeypass": "",
  "jksStorepass": "",
  "preCommand": "",
  "postCommand": "systemctl reload nginx"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| useSCP | bool | No | Use SCP instead of SFTP |
| format | string | Yes | Certificate format: `PEM`, `PFX`, or `JKS` |
| keyPath | string | Yes (PEM) | Private key file path on remote server |
| certPath | string | Yes | Certificate file path on remote server |
| certPathForServerOnly | string | No | Server-only certificate path (PEM) |
| certPathForIntermediaOnly | string | No | Intermediate certificate path (PEM) |
| pfxPassword | string | Yes (PFX) | Password for PFX format |
| jksAlias | string | Yes (JKS) | JKS alias |
| jksKeypass | string | Yes (JKS) | JKS key password |
| jksStorepass | string | Yes (JKS) | JKS store password |
| preCommand | string | No | Command to run before upload |
| postCommand | string | No | Command to run after upload |

**Template Variables for Commands:**
- `${CERTIMATE_DEPLOYER_CMDVAR_CERTIFICATE_PATH}` - Certificate file path
- `${CERTIMATE_DEPLOYER_CMDVAR_PRIVATEKEY_PATH}` - Private key file path
- `${CERTIMATE_DEPLOYER_CMDVAR_CERTIFICATE_SERVER_PATH}` - Server certificate path
- `${CERTIMATE_DEPLOYER_CMDVAR_CERTIFICATE_INTERMEDIA_PATH}` - Intermediate certificate path

---

## Kubernetes (k8s)

### k8s-secret (Kubernetes Secret)

```json
{
  "namespace": "default",
  "secretName": "tls-secret",
  "secretDataKeyForCrt": "tls.crt",
  "secretDataKeyForKey": "tls.key"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| namespace | string | Yes | Kubernetes namespace |
| secretName | string | Yes | Secret name |
| secretDataKeyForCrt | string | No | Key for certificate (default: tls.crt) |
| secretDataKeyForKey | string | No | Key for private key (default: tls.key) |

---

## Other Providers

### 1panel
```json
{ "websiteId": "1" }
```

### 1panel-console
```json
{ "websiteId": "1" }
```

### apisix (Apache APISIX)
```json
{
  "serverUrl": "http://apisix-admin:9180",
  "sslId": "ssl-xxx"
}
```

### azure-keyvault
```json
{
  "vaultName": "my-vault",
  "certificateName": "my-certificate"
}
```

### baishan-cdn
```json
{ "domain": "cdn.example.com" }
```

### baotapanel
```json
{ "siteId": "1" }
```

### baotapanel-console
```json
{ "siteId": "1" }
```

### baotapanelgo
```json
{ "siteId": "1" }
```

### baotapanelgo-console
```json
{ "siteId": "1" }
```

### baotawaf
```json
{ "siteId": "1" }
```

### baotawaf-console
```json
{ "siteId": "1" }
```

### bunny-cdn
```json
{
  "pullZoneId": "12345",
  "hostname": "cdn.example.com"
}
```

### byteplus-cdn
```json
{ "domain": "cdn.example.com" }
```

### cachefly
```json
{ "hostname": "cdn.example.com" }
```

### cdnfly
```json
{
  "domain": "cdn.example.com",
  "resourceId": "res-xxx"
}
```

### cpanel
```json
{
  "hostname": "cpanel.example.com",
  "domain": "example.com"
}
```

### dogecloud-cdn
```json
{ "domain": "cdn.example.com" }
```

### dokploy
```json
{ "appName": "my-app" }
```

### flexcdn
```json
{ "domain": "cdn.example.com" }
```

### flyio
```json
{ "appName": "my-app" }
```

### gcore-cdn
```json
{
  "resourceId": "res-xxx",
  "domain": "cdn.example.com"
}
```

### goedge
```json
{
  "clusterId": "1",
  "certificateId": "1"
}
```

### kong
```json
{
  "serverUrl": "http://kong-admin:8001",
  "certificateId": "cert-xxx"
}
```

### ksyun-cdn
```json
{ "domain": "cdn.example.com" }
```

### lecdn
```json
{ "domain": "cdn.example.com" }
```

### local (Local File)
```json
{
  "outputPath": "/etc/ssl/certimate",
  "format": "PEM",
  "pfxPassword": "",
  "jksAlias": "",
  "jksKeypass": "",
  "jksStorepass": ""
}
```

### mohua-mvh
```json
{ "domain": "example.com" }
```

### netlify
```json
{
  "siteId": "site-xxx",
  "domain": "example.com"
}
```

### nginxproxymanager
```json
{
  "serverUrl": "http://npm:81",
  "proxyHostId": "1"
}
```

### proxmoxve
```json
{
  "nodeName": "pve1",
  "realm": "pam"
}
```

### rainyun-rcdn
```json
{ "domain": "cdn.example.com" }
```

### rainyun-sslcenter
```json
{ "certificateId": "1" }
```

### ratpanel
```json
{ "siteId": "1" }
```

### ratpanel-console
```json
{ "siteId": "1" }
```

### s3 (S3 Compatible Storage)
```json
{
  "region": "us-east-1",
  "bucket": "my-bucket",
  "path": "/certificates/"
}
```

### safeline
```json
{
  "serverUrl": "https://safeline:9443",
  "resourceId": "1"
}
```

### synologydsm
```json
{ "domain": "example.com" }
```

### unicloud-webhost
```json
{
  "spaceId": "space-xxx",
  "domain": "example.com"
}
```

### upyun-cdn
```json
{ "domain": "cdn.example.com" }
```

### upyun-file
```json
{
  "bucket": "my-bucket",
  "path": "/certificates/"
}
```

### wangsu-cdn
```json
{ "domain": "cdn.example.com" }
```

### wangsu-cdnpro
```json
{ "domain": "cdn.example.com" }
```

### wangsu-certificate
```json
{ "domain": "example.com" }
```

### webhook
```json
{
  "url": "https://example.com/webhook",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer token"
  }
}
```

---

## CTCC Cloud (ctcccloud)

### ctcccloud-ao
```json
{
  "region": "sh",
  "instanceId": "instance-xxx"
}
```

### ctcccloud-cdn
```json
{ "domain": "cdn.example.com" }
```

### ctcccloud-cms
```json
{
  "region": "sh",
  "instanceId": "instance-xxx"
}
```

### ctcccloud-elb
```json
{
  "region": "sh",
  "loadbalancerId": "lb-xxx",
  "listenerId": "lsr-xxx"
}
```

### ctcccloud-faas
```json
{
  "region": "sh",
  "instanceId": "instance-xxx"
}
```

### ctcccloud-icdn
```json
{ "domain": "cdn.example.com" }
```

### ctcccloud-ldvn
```json
{ "domain": "cdn.example.com" }
```

---

## Domain Match Patterns

For CDN and load balancer providers that support domain matching:

| Pattern | Description | Example |
|---------|-------------|---------|
| `exact` | Exact domain match | `example.com` matches only `example.com` |
| `wildcard` | Wildcard domain match | `*.example.com` matches all subdomains |
| `certsan` | Match against certificate SANs | Uses domains from certificate |

---

## Quick Reference by Use Case

### CDN Deployment
- Aliyun: `aliyun-cdn`
- Tencent: `tencentcloud-cdn`
- Huawei: `huaweicloud-cdn`
- VolcEngine: `volcengine-cdn`
- Qiniu: `qiniu-cdn`
- Baidu: `baiducloud-cdn`

### Load Balancer Deployment
- Aliyun ALB: `aliyun-alb`
- Aliyun CLB: `aliyun-clb`
- Tencent CLB: `tencentcloud-clb`
- Huawei ELB: `huaweicloud-elb`

### Object Storage Deployment
- Aliyun OSS: `aliyun-oss`
- Tencent COS: `tencentcloud-cos`
- Huawei OBS: `huaweicloud-obs`
- AWS S3: `s3`
- VolcEngine TOS: `volcengine-tos`

### Server Deployment
- SSH: `ssh`
- Local: `local`

### Container Deployment
- Kubernetes: `k8s-secret`

### Panel Deployment
- 1Panel: `1panel`
- BaoTa Panel: `baotapanel`
- Dokploy: `dokploy`
