# lynx-mcp-server

Golang server that implements the MCP (Model Control Protocol) for interacting with Lynx Reservations system.

## Prerequisites

- Go 1.23.2 or later
- Docker

## Environment Variables

The following environment variables are required:

- `LYNX_USERNAME`: Your Lynx Reservations username
- `LYNX_PASSWORD`: Your Lynx Reservations password  
- `LYNX_COMPANY_CODE`: Your Lynx company code

Use `.env` file to work locally.

## How to build

```sh
docker build -t lynx-mcp-server .  
```

### How to run

```sh
docker run --env-file .env -p 9600:9600 lynx-mcp-server  
```

## How to test

```sh
go run ./client
```