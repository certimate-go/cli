# Provider Configuration Reference

This directory contains configuration reference documents for all supported providers in Certimate.

## Quick Reference

### DNS Providers

| Provider | Config Fields | Required |
|----------|---------------|----------|
| 1Panel | `serverUrl`, `apiVersion`, `apiKey`, `allowInsecureConnections` | `serverUrl`, `apiKey` |
| 35cn | `username`, `apiPassword` | `username`, `apiPassword` |
| 51DNScom | `apiKey`, `apiSecret` | `apiKey`, `apiSecret` |
| Aliyun | `accessKeyId`, `accessKeySecret`, `resourceGroupId` | `accessKeyId`, `accessKeySecret` |
| AWS | `accessKeyId`, `secretAccessKey` | `accessKeyId`, `secretAccessKey` |
| Azure | `tenantId`, `clientId`, `clientSecret`, `subscriptionId`, `resourceGroupName`, `cloudName` | `tenantId`, `clientId`, `clientSecret` |
| BaiduCloud | `accessKeyId`, `secretAccessKey` | `accessKeyId`, `secretAccessKey` |
| Cloudflare | `dnsApiToken`, `zoneApiToken` | `dnsApiToken` |
| ClouDNS | `authId`, `authPassword` | `authId`, `authPassword` |
| Constellix | `apiKey`, `secretKey` | `apiKey`, `secretKey` |
| DeSEC | `token` | `token` |
| DigitalOcean | `accessToken` | `accessToken` |
| DNSExit | `apiKey` | `apiKey` |
| DNSLA | `apiId`, `apiSecret` | `apiId`, `apiSecret` |
| DNSMadeEasy | `apiKey`, `apiSecret` | `apiKey`, `apiSecret` |
| DuckDNS | `token` | `token` |
| Dynu | `apiKey` | `apiKey` |
| Dynv6 | `httpToken` | `httpToken` |
| GoDaddy | `apiKey`, `apiSecret` | `apiKey`, `apiSecret` |
| HuaweiCloud | `accessKeyId`, `secretAccessKey`, `enterpriseProjectId` | `accessKeyId`, `secretAccessKey` |
| Infomaniak | `accessToken` | `accessToken` |
| IONOS | `apiKeyPublicPrefix`, `apiKeySecret` | `apiKeyPublicPrefix`, `apiKeySecret` |
| JDCloud | `accessKeyId`, `accessKeySecret` | `accessKeyId`, `accessKeySecret` |
| Ksyun | `accessKeyId`, `secretAccessKey` | `accessKeyId`, `secretAccessKey` |
| Linode | `accessToken` | `accessToken` |
| Namecheap | `username`, `apiKey` | `username`, `apiKey` |
| NameDotCom | `username`, `apiToken` | `username`, `apiToken` |
| NameSilo | `apiKey` | `apiKey` |
| Netcup | `customerNumber`, `apiKey`, `apiPassword` | `customerNumber`, `apiKey`, `apiPassword` |
| Netlify | `apiToken` | `apiToken` |
| NS1 | `apiKey` | `apiKey` |
| OVHcloud | `endpoint`, `authMethod`, `applicationKey`, `applicationSecret`, `consumerKey`, `clientId`, `clientSecret` | `endpoint`, `authMethod` |
| Porkbun | `apiKey`, `secretApiKey` | `apiKey`, `secretApiKey` |
| PowerDNS | `serverUrl`, `apiKey` | `serverUrl`, `apiKey` |
| QingCloud | `accessKeyId`, `secretAccessKey` | `accessKeyId`, `secretAccessKey` |
| Qiniu | `accessKey`, `secretKey` | `accessKey`, `secretKey` |
| RainYun | `apiKey` | `apiKey` |
| RFC2136 | `host`, `port`, `tsigAlgorithm`, `tsigKey`, `tsigSecret` | `host`, `port` |
| Spaceship | `apiKey`, `apiSecret` | `apiKey`, `apiSecret` |
| TencentCloud | `secretId`, `secretKey` | `secretId`, `secretKey` |
| TodayNIC | `userId`, `apiKey` | `userId`, `apiKey` |
| UCloud | `privateKey`, `publicKey`, `projectId` | `privateKey`, `publicKey` |
| Vercel | `apiAccessToken`, `teamId` | `apiAccessToken` |
| VolcEngine | `accessKeyId`, `secretAccessKey` | `accessKeyId`, `secretAccessKey` |
| Vultr | `apiKey` | `apiKey` |
| Wangsu | `accessKeyId`, `accessKeySecret`, `apiKey` | `accessKeyId`, `accessKeySecret` |
| Westcn | `username`, `apiPassword` | `username`, `apiPassword` |
| Xinnet | `agentId`, `apiPassword` | `agentId`, `apiPassword` |

### CDN/Deployment Providers

| Provider | Config Fields | Required |
|----------|---------------|----------|
| Akamai | `host`, `clientToken`, `clientSecret`, `accessToken` | `host`, `clientToken`, `clientSecret`, `accessToken` |
| APISIX | `serverUrl`, `apiKey`, `allowInsecureConnections` | `serverUrl`, `apiKey` |
| ArvanCloud | `apiKey` | `apiKey` |
| Baishan | `apiToken` | `apiToken` |
| BaotaPanel | `serverUrl`, `apiKey`, `allowInsecureConnections` | `serverUrl`, `apiKey` |
| Bunny | `apiKey` | `apiKey` |
| BytePlus | `accessKey`, `secretKey` | `accessKey`, `secretKey` |
| CacheFly | `apiToken` | `apiToken` |
| Cdnfly | `serverUrl`, `apiKey`, `apiSecret`, `allowInsecureConnections` | `serverUrl`, `apiKey`, `apiSecret` |
| FlexCDN | `serverUrl`, `apiRole`, `accessKeyId`, `accessKey`, `allowInsecureConnections` | `serverUrl`, `apiRole`, `accessKeyId`, `accessKey` |
| FlyIO | `apiToken` | `apiToken` |
| Gcore | `apiToken` | `apiToken` |
| GoEdge | `serverUrl`, `apiRole`, `accessKeyId`, `accessKey`, `allowInsecureConnections` | `serverUrl`, `apiRole`, `accessKeyId`, `accessKey` |
| Hetzner | `apiToken` | `apiToken` |
| Hostingde | `apiKey` | `apiKey` |
| Hostinger | `apiToken` | `apiToken` |
| Kong | `serverUrl`, `apiToken`, `allowInsecureConnections` | `serverUrl` |
| Kubernetes | `kubeConfig` | Optional |
| LeCDN | `serverUrl`, `apiVersion`, `apiRole`, `username`, `password`, `allowInsecureConnections` | `serverUrl`, `apiVersion`, `apiRole`, `username`, `password` |
| Mattermost | `serverUrl`, `username`, `password`, `channelId` | `serverUrl`, `username`, `password` |
| Mohua | `username`, `apiPassword` | `username`, `apiPassword` |
| NginxProxyManager | `serverUrl`, `authMethod`, `username`, `password`, `apiToken`, `allowInsecureConnections` | `serverUrl`, `authMethod` |
| ProxmoxVE | `serverUrl`, `apiToken`, `apiTokenSecret`, `allowInsecureConnections` | `serverUrl`, `apiToken` |
| RatPanel | `serverUrl`, `accessTokenId`, `accessToken`, `allowInsecureConnections` | `serverUrl`, `accessTokenId`, `accessToken` |
| S3 | `endpoint`, `accessKey`, `secretKey`, `signatureVersion`, `usePathStyle`, `allowInsecureConnections` | `endpoint`, `accessKey`, `secretKey` |
| SafeLine | `serverUrl`, `apiToken`, `allowInsecureConnections` | `serverUrl`, `apiToken` |
| SSH | `host`, `port`, `authMethod`, `username`, `password`, `key`, `keyPassphrase`, `jumpServers` | `host`, `port`, `authMethod`, `username` |
| SynologyDSM | `serverUrl`, `username`, `password`, `totpSecret`, `allowInsecureConnections` | `serverUrl`, `username`, `password` |
| TechnitiumDNS | `serverUrl`, `apiToken`, `allowInsecureConnections` | `serverUrl`, `apiToken` |
| UniCloud | `username`, `password` | `username`, `password` |
| Upyun | `username`, `password` | `username`, `password` |

### Notification Providers

| Provider | Config Fields | Required |
|----------|---------------|----------|
| DingTalkBot | `webhookUrl`, `secret`, `customPayload` | `webhookUrl` |
| DiscordBot | `botToken`, `channelId` | `botToken` |
| Email | `smtpHost`, `smtpPort`, `smtpTls`, `username`, `password`, `senderAddress`, `senderName`, `receiverAddress`, `allowInsecureConnections` | `smtpHost`, `smtpPort`, `smtpTls`, `username`, `password`, `senderAddress`, `senderName` |
| LarkBot | `webhookUrl`, `secret`, `customPayload` | `webhookUrl` |
| SlackBot | `botToken`, `channelId` | `botToken` |
| TelegramBot | `botToken`, `chatId` | `botToken` |
| Webhook | `url`, `method`, `headers`, `data`, `allowInsecureConnections` | `url` |
| WeComBot | `webhookUrl`, `customPayload` | `webhookUrl` |

### ACME Providers

| Provider | Config Fields | Required |
|----------|---------------|----------|
| ACMECA | `endpoint`, `eabKid`, `eabHmacKey` | `endpoint` |
| ACMEDNS | `serverUrl`, `credentials` | `serverUrl`, `credentials` |
| ACMEHttpReq | `endpoint`, `mode`, `username`, `password` | `endpoint` |
| ActalisSSL | `eabKid`, `eabHmacKey` | Optional |
| GlobalSignAtlas | `eabKid`, `eabHmacKey` | Optional |
| GoogleTrustServices | `eabKid`, `eabHmacKey` | Optional |
| LiteSSL | `eabKid`, `eabHmacKey` | Optional |
| SSLCom | `eabKid`, `eabHmacKey` | Optional |
| ZeroSSL | `eabKid`, `eabHmacKey` | Optional |

## Configuration Examples

### Cloudflare DNS API Token

```json
{
  "dnsApiToken": "your-cloudflare-dns-api-token"
}
```

### Aliyun Access Key

```json
{
  "accessKeyId": "your-access-key-id",
  "accessKeySecret": "your-access-key-secret"
}
```

### AWS Route53

```json
{
  "accessKeyId": "AKIAIOSFODNN7EXAMPLE",
  "secretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

### Tencent Cloud

```json
{
  "secretId": "your-secret-id",
  "secretKey": "your-secret-key"
}
```

## Sensitive Fields

The following fields are automatically masked in output (shown as `********`):

- `secretAccessKey`
- `secretKey`
- `accessKeySecret`
- `apiSecret`
- `apiKey`
- `password`
- `token`
- `privateKey`
- `secret`
- `authorizationToken`
- `secretId`

Use the `--reveal` flag to show actual values (use with caution).

## Related Documentation

- [Main SKILL.md](../SKILL.md) - Skill usage instructions
- [Certimate Provider Documentation](https://deepwiki.com/certimate-go/certimate)
