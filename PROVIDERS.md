# DNS Provider Setup Guide

This guide provides detailed setup instructions, API permissions, and usage examples for each supported DNS provider.

## Table of Contents

- [AWS Route53](#aws-route53)
- [Cloudflare DNS](#cloudflare-dns)
- [Azure DNS](#azure-dns)
- [DigitalOcean DNS](#digitalocean-dns)
- [OVH DNS](#ovh-dns)

---

## AWS Route53

### Setup Requirements

1. **Create IAM User:**
   - Go to AWS IAM Console
   - Create new user with programmatic access
   - Attach policy: `AmazonRoute53FullAccess` (or create custom policy)

2. **Required Permissions:**
   ```json
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Action": [
           "route53:ListHostedZones",
           "route53:ChangeResourceRecordSets"
         ],
         "Resource": "*"
       }
     ]
   }
   ```

3. **Get Zone ID:**
   - Go to Route53 console
   - Find your hosted zone
   - Copy the Zone ID (format: `Z1D633PJN98FT9`)

### Authentication

| Field | Value |
|-------|-------|
| **Username** | AWS Access Key ID |
| **Password** | AWS Secret Access Key |

### URL Format

```
http://localhost:8080/aws/[ZONE_ID]/?ip=[IP_ADDRESS]&hostname=[HOSTNAME]
```

### Usage Examples

```bash
# Update A record for test.example.com
curl -u "AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  "http://localhost:8080/aws/Z1D633PJN98FT9/?ip=192.168.1.100&hostname=test.example.com"

# Update root domain
curl -u "AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  "http://localhost:8080/aws/Z1D633PJN98FT9/?ip=192.168.1.100&hostname=example.com"
```

---

## Cloudflare DNS

### Setup Requirements

1. **Get API Token:**
   - Go to Cloudflare Dashboard → My Profile → API Tokens
   - Create token with `Zone:Edit` permissions
   - Select specific zones or all zones

2. **Required Permissions:**
   - Zone: Zone Settings: Read
   - Zone: Zone: Read  
   - Zone: DNS: Edit

3. **Get Zone Name:**
   - Use your domain name (e.g., `example.com`)

### Authentication

| Field | Value |
|-------|-------|
| **Username** | Zone Name (domain) |
| **Password** | API Token |

### URL Format

```
http://localhost:8080/cloudflare/?ip=[IP_ADDRESS]&hostname=[HOSTNAME]
```

### Usage Examples

```bash
# Update subdomain
curl -u "example.com:your-cloudflare-api-token" \
  "http://localhost:8080/cloudflare/?ip=192.168.1.100&hostname=test.example.com"

# Update root domain  
curl -u "example.com:your-cloudflare-api-token" \
  "http://localhost:8080/cloudflare/?ip=192.168.1.100&hostname=example.com"
```

---

## Azure DNS

### Setup Requirements

1. **Create App Registration:**
   - Go to Azure Portal → Azure Active Directory → App registrations
   - Create new registration
   - Note the Application (client) ID

2. **Create Client Secret:**
   - In your app registration → Certificates & secrets
   - Create new client secret
   - Copy the secret value

3. **Assign Permissions:**
   - Go to your DNS Zone → Access control (IAM)
   - Add role assignment: `DNS Zone Contributor`
   - Assign to your application

4. **Get Required IDs:**
   - **Tenant ID**: Azure AD → Properties → Tenant ID
   - **Subscription ID**: Subscriptions → Your subscription → Subscription ID
   - **Resource Group**: Resource groups → Your DNS zone's resource group
   - **Zone Name**: Your domain name

### Authentication

| Field | Value |
|-------|-------|
| **Username** | Application (Client) ID |
| **Password** | Client Secret |

### URL Format

```
http://localhost:8080/azure/[TENANT_ID]/[SUBSCRIPTION_ID]/[RESOURCE_GROUP]/[ZONE_NAME]/?ip=[IP_ADDRESS]&hostname=[HOSTNAME]
```

### Usage Examples

```bash
# Update subdomain
curl -u "12345678-1234-1234-1234-123456789012:your-client-secret" \
  "http://localhost:8080/azure/tenant-id/subscription-id/my-resource-group/example.com/?ip=192.168.1.100&hostname=test.example.com"

# Update root domain (uses @ record)
curl -u "12345678-1234-1234-1234-123456789012:your-client-secret" \
  "http://localhost:8080/azure/tenant-id/subscription-id/my-resource-group/example.com/?ip=192.168.1.100&hostname=example.com"
```

---

## DigitalOcean DNS

### Setup Requirements

1. **Create API Token:**
   - Go to DigitalOcean Control Panel → API
   - Generate new token with read/write scope

2. **Add Domain:**
   - Go to Networking → Domains
   - Add your domain to DigitalOcean DNS

3. **Required Permissions:**
   - Domain records read/write access

### Authentication

| Field | Value |
|-------|-------|
| **Username** | Domain Name |
| **Password** | API Token |

### URL Format

```
http://localhost:8080/digitalocean/?ip=[IP_ADDRESS]&hostname=[HOSTNAME]
```

### Usage Examples

```bash
# Update subdomain
curl -u "example.com:your-digitalocean-api-token" \
  "http://localhost:8080/digitalocean/?ip=192.168.1.100&hostname=test.example.com"

# Update root domain (@ record)
curl -u "example.com:your-digitalocean-api-token" \
  "http://localhost:8080/digitalocean/?ip=192.168.1.100&hostname=example.com"
```

---

## OVH DNS

### Setup Requirements

1. **Create Application:**
   - Go to OVH API Console: https://eu.api.ovh.com/createToken/
   - Create application to get Application Key and Application Secret

2. **Generate Consumer Key:**
   - Use the API console to generate Consumer Key
   - **Required API Permissions:**
     ```
     GET /domain/zone/*/record
     POST /domain/zone/*/record
     PUT /domain/zone/*/record/*
     POST /domain/zone/*/refresh
     ```

3. **Choose Endpoint:**
   - **EU**: `eu` (Europe) → https://eu.api.ovh.com/1.0
   - **CA**: `ca` (Canada) → https://ca.api.ovh.com/1.0  
   - **US**: `us` (United States) → https://api.ovh.com/1.0
   - **AU**: `au` (Australia) → https://au.api.ovh.com/1.0

### Authentication

| Field | Value |
|-------|-------|
| **Username** | Application Secret |
| **Password** | Consumer Key |

### URL Format

```
http://localhost:8080/ovh/[ENDPOINT]/[DOMAIN]/[APPLICATION_KEY]/?ip=[IP_ADDRESS]&hostname=[HOSTNAME]
```

### Usage Examples

```bash
# Europe endpoint
curl -u "application-secret:consumer-key" \
  "http://localhost:8080/ovh/eu/example.com/your-app-key/?ip=192.168.1.100&hostname=test.example.com"

# US endpoint  
curl -u "application-secret:consumer-key" \
  "http://localhost:8080/ovh/us/example.com/your-app-key/?ip=192.168.1.100&hostname=test.example.com"

# Canada endpoint
curl -u "application-secret:consumer-key" \
  "http://localhost:8080/ovh/ca/example.com/your-app-key/?ip=192.168.1.100&hostname=test.example.com"

# Australia endpoint
curl -u "application-secret:consumer-key" \
  "http://localhost:8080/ovh/au/example.com/your-app-key/?ip=192.168.1.100&hostname=test.example.com"
```

### Troubleshooting OVH

**Common Issues:**

1. **"This application key is invalid"**
   - Verify Application Key is correct
   - Check you're using the right endpoint region

2. **"This call has not been granted"**
   - Consumer Key lacks required permissions
   - Regenerate Consumer Key with proper API rights

3. **Records created but not visible**
   - Zone refresh may have failed
   - Check DNS propagation (can take up to 24 hours)

**Testing API Access:**
```bash
# Test if your credentials work
curl -X GET "https://eu.api.ovh.com/1.0/domain/zone" \
  -H "X-Ovh-Application: YOUR_APP_KEY" \
  -H "X-Ovh-Consumer: YOUR_CONSUMER_KEY"
```

---

## General Usage Notes

### IP Address Validation
- All providers validate IP addresses before making API calls
- Only IPv4 addresses are currently supported
- Invalid IPs return `400 Bad Request`

### Hostname Requirements  
- Must be valid FQDN or match domain/zone
- Subdomain extraction is automatic
- Root domain updates use appropriate record types (@, blank, etc.)

### Error Handling
- All providers return `200 OK` on success
- Error responses include descriptive messages
- Check application logs for detailed error information

### Rate Limiting
- Respect provider API rate limits
- Most providers have generous limits for DNS operations
- Consider caching for high-frequency updates

### SSL/TLS
- Application runs on HTTP by default (localhost only)
- Use reverse proxy (nginx, Apache) for HTTPS in production
- Never expose credentials over unencrypted connections