# Cloud DDNS

This is a Go application that acts as a Dynamic DNS (DDNS) bridge, converting DynDNS-formatted HTTP requests into API calls for cloud DNS services. Here's what it does:

## Overview

The application provides a unified interface for updating DNS records across different cloud providers (AWS Route53 and Cloudflare) using the standard DynDNS HTTP protocol format.

## Architecture

**Main Components:**
- **HTTP Server** (`main.go`, `http.go`) - Handles incoming requests with basic authentication
- **AWS Route53 Handler** (`aws.go`) - Manages DNS updates for AWS Route53 
- **Cloudflare Handler** (`cloudflare.go`) - Manages DNS updates for Cloudflare DNS
- **Digital Ocean Handler** (`digitalocean.go`) - Maanages DNS updates for Ditital Ocean

## How It Works

**Request Format:**
- AWS: `http://server:port/aws/[zoneid]?ip=[ip_address]&hostname=[hostname]`
- Cloudflare: `http://server:port/cloudflare/?ip=[ip_address]&hostname=[hostname]`
- http://server:port/digitalocean/?ip=[ip_address]&hostname=[hostname]

**Authentication:**
- Uses HTTP Basic Authentication
- Username/password serve different purposes for each provider:
  - AWS: username = access key, password = secret key
  - Cloudflare: username = zone name, password = API token
  - Digital Ocean username = zone name, password = API Token

**Process Flow:**
1. Client sends DynDNS-formatted request
2. Server validates credentials and form data
3. Extracts IP address and hostname from request
4. Calls appropriate cloud provider API
5. Creates or updates DNS A record
6. Returns success/failure response

## Key Features

- **Dual Provider Support** - Works with both AWS Route53 and Cloudflare
- **Standard DynDNS Protocol** - Compatible with existing DynDNS clients
- **Automatic Record Management** - Creates new records or updates existing ones
- **Security** - Requires authentication and validates input
- **Logging** - Uses syslog for operation logging
- **Configurable** - Supports custom IP/port binding

This is particularly useful for home networks or small businesses that want to use enterprise cloud DNS services with standard DynDNS-compatible routers or clients.