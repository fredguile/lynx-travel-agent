# Docker Compose Setup

This directory contains Docker Compose configurations for running the Pan PAC
helper services with automated workflows and integrations.

## Overview

The Docker Compose setup provides a complete environment for running:

- **n8n** - Workflow automation platform
- **PostgreSQL** - Database for n8n workflows
- **Traefik** - Reverse-proxy with automatic SSL certificates
- **Lynx MCP Server** - MCP Server for Lynx Reservations integration

## Configuration Files

### `compose.yaml` (Main Configuration)

Full-featured setup with PostgreSQL database backend for production use.

**Services:**

- **Traefik** - Reverse-proxy with automatic SSL certificate management
- **PostgreSQL 16** - Production database with optimized settings
- **n8n** - Workflow automation with PostgreSQL backend
- **Lynx MCP Server** - MCP Server for Lynx Reservations

### `compose-sqlite.yaml` (Lightweight Configuration - NOT IN USE)

Simplified setup using SQLite for development or lightweight deployments.

**Services:**

- **Traefik** - Reverse proxy with automatic SSL certificate management
- **n8n** - Workflow automation with SQLite backend

## How to deploy

### 1. Create Required Volumes

```bash
docker volume create db_data
docker volume create n8n_data
docker volume create traefik_data
```

### 2. Set Environment Variables

Create a `.env` file or use the provided `.envrc` as a reference. Required
variables:

**For PostgreSQL setup:**

```bash
POSTGRES_PASSWORD=your_secure_password
POSTGRES_USER=postgres
POSTGRES_DB=postgres
POSTGRES_NON_ROOT_USER=your_user
POSTGRES_NON_ROOT_PASSWORD=your_user_password
```

**For Traefik SSL:**

```bash
SSL_EMAIL=your-email@domain.com
DOMAIN_NAME=your-domain.com
N8N_SUBDOMAIN=n8n
LYNX_MCP_SUBDOMAIN=mcp
```

**For Lynx MCP Server:**

```bash
LYNX_USERNAME=your_lynx_username
LYNX_PASSWORD=your_lynx_password
LYNX_COMPANY_CODE=your_company_code
```

### 3. Start Services

```bash
docker-compose up -d
```

## Service Details

### Traefik

- **Ports:** 80, 443
- **Features:** Automatic SSL certificates via Let's Encrypt

### PostgreSQL (compose.yaml only)

- **Port:** 5432 (internal)
- **Optimizations:** Configured with production-ready settings
- **Health Check:** Automatic health monitoring

### n8n

- **Port:** 5678 (internal, accessible via Traefik)
- **Access:** `https://n8n.your-domain.com` (or configured subdomain)
- **Features:**
  - Production mode with optimized settings
  - Internal runners for better performance
  - File system access via `./local-files` volume

### Lynx MCP Server

- **Port:** 9600 (internal, accessible via Traefik)
- **Access:** `https://mcp.your-domain.com` (or configured subdomain)
- **Purpose:** MCP Server for Lynx Reservations integration

## Volumes

- `db_data` - PostgreSQL database storage
- `n8n_data` - n8n workflow and configuration storage
- `traefik_data` - SSL certificates and Traefik configuration
- `./local-files` - Local file system access for n8n workflows

## Security Features

- **SSL/TLS:** Automatic certificate generation and renewal
- **Security Headers:** HSTS, XSS protection, content type sniffing prevention
- **Network Isolation:** Services only exposed through Traefik
- **Health Checks:** Database health monitoring

## Monitoring and Logs

```bash
# View all service logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f n8n
docker-compose logs -f postgres
docker-compose logs -f traefik

# Check service status
docker-compose ps
```

## Troubleshooting

### SSL Certificate Issues

- Ensure `SSL_EMAIL` is set correctly
- Check that ports 80 and 443 are accessible
- Verify domain DNS is pointing to your server

### Database Connection Issues

- Ensure PostgreSQL volume is created
- Check environment variables are set correctly
- Verify health check status: `docker-compose ps postgres`

### n8n Access Issues

- Check Traefik logs for routing issues
- Verify subdomain configuration
- Ensure all required environment variables are set

## Backup and Recovery

### Database Backup

```bash
# PostgreSQL backup
docker-compose exec postgres pg_dump -U postgres postgres > backup.sql

# n8n data backup
docker run --rm -v n8n_data:/data -v $(pwd):/backup alpine tar czf /backup/n8n_backup.tar.gz -C /data .
```

### Restore

```bash
# PostgreSQL restore
docker-compose exec -T postgres psql -U postgres postgres < backup.sql

# n8n data restore
docker run --rm -v n8n_data:/data -v $(pwd):/backup alpine tar xzf /backup/n8n_backup.tar.gz -C /data
```
