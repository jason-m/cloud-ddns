# Cloud DDNS

This is a Go application that acts as a Dynamic DNS (DDNS) bridge, converting DynDNS-formatted HTTP requests into API calls for cloud DNS services.

## Overview

The application provides a unified interface for updating DNS records across different cloud providers using the standard DynDNS HTTP protocol format.

## Supported Providers

- **AWS Route53** - Amazon's DNS service
- **Cloudflare DNS** - Cloudflare's DNS service  
- **Azure DNS** - Microsoft Azure's DNS service
- **DigitalOcean DNS** - DigitalOcean's DNS service

## Architecture

**Main Components:**
- **HTTP Server** (`main.go`, `http.go`) - Handles incoming requests with basic authentication
- **AWS Route53 Handler** (`aws.go`) - Manages DNS updates for AWS Route53 
- **Cloudflare Handler** (`cloudflare.go`) - Manages DNS updates for Cloudflare DNS
- **Azure DNS Handler** (`azure.go`) - Manages DNS updates for Azure DNS
- **DigitalOcean Handler** (`digitalocean.go`) - Manages DNS updates for DigitalOcean DNS

## Request Formats

**AWS Route53:**
```
http://server:port/aws/[zoneid]/?ip=[ip_address]&hostname=[hostname]
```

**Cloudflare:**
```
http://server:port/cloudflare/?ip=[ip_address]&hostname=[hostname]
```

**Azure DNS:**
```
http://server:port/azure/[tenantid]/[subscriptionid]/[resource-group]/[zone-name]/?ip=[ip_address]&hostname=[hostname]
```

**DigitalOcean:**
```
http://server:port/digitalocean/?ip=[ip_address]&hostname=[hostname]
```

## Authentication

Uses HTTP Basic Authentication where username/password serve different purposes for each provider:

| Provider | Username | Password |
|----------|----------|----------|
| **AWS Route53** | Access Key | Secret Key |
| **Cloudflare** | Zone Name | API Token |
| **Azure DNS** | Client ID | Client Secret |
| **DigitalOcean** | Domain Name | API Token |

## Setup Instructions

### AWS Route53
1. Create IAM user with Route53 permissions
2. Generate Access Key and Secret Key
3. Use Zone ID from Route53 console

### Cloudflare
1. Get API Token from Cloudflare dashboard
2. Use zone name as username

### Azure DNS
1. Create App Registration in Azure Portal
2. Get Client ID, Client Secret, and Tenant ID
3. Assign DNS Zone Contributor role
4. Find your Subscription ID

### DigitalOcean
1. Generate API Token from DigitalOcean control panel
2. Use domain name as username

## Usage Examples

**AWS Route53:**
```bash
curl -u "AKIA....:wJalrXUtnFEMI...." \
  "http://localhost:8080/aws/Z1D633PJN98FT9/?ip=192.168.1.100&hostname=test.example.com"
```

**Cloudflare:**
```bash
curl -u "example.com:your-api-token" \
  "http://localhost:8080/cloudflare/?ip=192.168.1.100&hostname=test.example.com"
```

**Azure DNS:**
```bash
curl -u "client-id:client-secret" \
  "http://localhost:8080/azure/tenant-id/subscription-id/resource-group/zone-name/?ip=192.168.1.100&hostname=test.example.com"
```

**DigitalOcean:**
```bash
curl -u "example.com:your-api-token" \
  "http://localhost:8080/digitalocean/?ip=192.168.1.100&hostname=test.example.com"
```

## Process Flow

1. Client sends DynDNS-formatted request
2. Server validates credentials and form data
3. Extracts IP address and hostname from request
4. Calls appropriate cloud provider API
5. Creates or updates DNS A record
6. Returns success/failure response

## Key Features

- **Multi-Provider Support** - Works with AWS Route53, Cloudflare, Azure DNS, and DigitalOcean
- **Standard DynDNS Protocol** - Compatible with existing DynDNS clients
- **Automatic Record Management** - Creates new records or updates existing ones
- **Security** - Requires authentication and validates input
- **Logging** - Uses syslog for operation logging
- **Configurable** - Supports custom IP/port binding

## Build and run

1. Install Go dependencies:
```bash
go mod init cloud-ddns
go get
```

2. Build the application:
```bash
go build -o cloud-ddns
```

3. Run the application:
```bash
./cloud-ddns [ip] [port]
```

## Use Cases

This is particularly useful for:
- Home networks with dynamic IP addresses
- Small businesses using enterprise cloud DNS services
- IoT devices that need to update their DNS records
- Any scenario requiring DynDNS-compatible clients with cloud DNS providers