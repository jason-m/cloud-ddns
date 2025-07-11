# Cloud DDNS

A Go application that acts as a Dynamic DNS (DDNS) bridge, converting DynDNS-formatted HTTP requests into API calls for cloud DNS services.

## Overview

Cloud DDNS provides a unified interface for updating DNS records across different cloud providers using the standard DynDNS HTTP protocol format. This makes it compatible with existing DynDNS clients while leveraging modern cloud DNS services.

## Supported Providers

- **AWS Route53** - Amazon's DNS service
- **Cloudflare DNS** - Cloudflare's DNS service  
- **Azure DNS** - Microsoft Azure's DNS service
- **DigitalOcean DNS** - DigitalOcean's DNS service
- **OVH DNS** - OVH's DNS service (EU, CA, US, AU endpoints)

## Key Features

- **Multi-Provider Support** - Single application supporting 5 major cloud DNS providers
- **Standard DynDNS Protocol** - Compatible with existing DynDNS clients and routers
- **Automatic Record Management** - Creates new records or updates existing ones (UPSERT)
- **Security** - HTTP Basic Authentication for all providers
- **Logging** - Comprehensive syslog integration for monitoring
- **Configurable** - Custom IP/port binding support
- **Regional Support** - Multiple endpoints for global providers

## Quick Start

1. **Build the application:**
   ```bash
   go mod init cloud-ddns
   go get
   go build -o cloud-ddns
   ```

2. **Run the application:**
   ```bash
   ./cloud-ddns [ip] [port]
   # Example: ./cloud-ddns 127.0.0.1 8080
   ```

3. **Update DNS records:**
   ```bash
   curl -u "username:password" \
     "http://localhost:8080/provider/[params]/?ip=192.168.1.100&hostname=test.example.com"
   ```

## URL Formats

Each provider has a specific URL format and authentication method:

| Provider | URL Format |
|----------|------------|
| **AWS Route53** | `/aws/[zoneid]/?ip=x.x.x.x&hostname=host.domain.com` |
| **Cloudflare** | `/cloudflare/?ip=x.x.x.x&hostname=host.domain.com` |
| **Azure DNS** | `/azure/[tenantid]/[subscriptionid]/[resource-group]/[zone-name]/?ip=x.x.x.x&hostname=host.domain.com` |
| **DigitalOcean** | `/digitalocean/?ip=x.x.x.x&hostname=host.domain.com` |
| **OVH** | `/ovh/[endpoint]/[domain]/[appkey]/?ip=x.x.x.x&hostname=host.domain.com` |

## Authentication

All providers use HTTP Basic Authentication with provider-specific credentials:

| Provider | Username | Password |
|----------|----------|----------|
| **AWS Route53** | Access Key | Secret Key |
| **Cloudflare** | Zone Name | API Token |
| **Azure DNS** | Client ID | Client Secret |
| **DigitalOcean** | Domain Name | API Token |
| **OVH** | Application Secret | Consumer Key |

## Use Cases

- **Home Networks** - Dynamic IP addresses that need DNS updates
- **Small Businesses** - Enterprise cloud DNS with simple DDNS clients
- **IoT Devices** - Automated DNS record updates for connected devices
- **Legacy Systems** - DynDNS-compatible clients with modern cloud DNS
- **Multi-Cloud** - Single interface for multiple DNS providers

## Documentation

- **[providers.md](providers.md)** - Detailed setup instructions and examples for each provider
- **Configuration** - Command line options and deployment guidance
- **API Reference** - Complete endpoint documentation

## Architecture

The application consists of:
- **HTTP Server** - Handles incoming requests with basic authentication
- **Provider Handlers** - Individual modules for each DNS service
- **Common Functions** - Shared validation and logging functionality

## Security Notes

- Runs on localhost (127.0.0.1) by default for security
- Use a reverse proxy with SSL for external access
- All API credentials are passed via HTTP Basic Auth
- Input validation on IP addresses and hostnames
- Comprehensive logging for audit trails

## Contributing

When adding new DNS providers:
1. Follow the existing handler pattern
2. Implement provider-specific authentication
3. Support both create and update operations (UPSERT)
4. Add comprehensive error handling and logging
5. Update documentation in providers.md

## License

This project is open source. See LICENSE file for details.