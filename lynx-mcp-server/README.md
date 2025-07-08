# lynx-mcp-server

Golang server that implements the MCP (Model Control Protocol) for interacting with Lynx Reservations system.

Latest docker image is automatically published to https://ghcr.io/fredguile/lynx-mcp-server.

## Prerequisites

- Go 1.23.10 or later
- Docker

## Environment Variables

The following environment variables are required:

- `LYNX_USERNAME`: Your Lynx Reservations username
- `LYNX_PASSWORD`: Your Lynx Reservations password  
- `LYNX_COMPANY_CODE`: Your Lynx company code
- `BEARER_TOKEN`: Arbitrary JWT token to secure the server

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
go run ./cmd/lynxmcpserver.go

go run ./cmd/lynxmcpclient.go --command 'file_search_by_party_name --partyName=LASTNAME'

go run ./cmd/lynxmcpclient.go --command 'file_search_by_file_reference --fileReference=FTXXXXXXXXX'

go run ./cmd/lynxmcpclient.go --command 'retrieve_itinerary --fileIdentifier=XXX'

go run ./cmd/lynxmcpclient.go --command 'file_documents_by_transaction_reference --fileIdentifier=XXX --transactionIdentifier=XXX'

go run ./cmd/lynxmcpclient.go --command 'attachment_upload --binary --identifier=YYY --fileName=attachment.pdf'

go run ./cmd/lynxmcpclient.go --command 'file_document_save_details --fileIdentifier=XXX --transactionIdentifier=XXX --name=document --content "<span>test</span>" --type=SUPP --attachmentUrl=/documents/file/f16476987/d20250709064401.pdf'
```