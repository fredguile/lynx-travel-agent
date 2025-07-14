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

go run ./cmd/lynxmcpclient.go --command 'file_document_save_details --fileIdentifier=XXX --name=document --content "<span>test</span>" --type=SUPP --attachmentUrl=/documents/file/f16476987/d20250709064401.pdf'

go run ./cmd/lynxmcpclient.go --command 'transaction_document_save_details --fileIdentifier=XXX --transactionIdentifier=XXX --name=document --content "<span>test</span>" --type=SUPP --attachmentUrl=/documents/file/f16476987/d20250709064401.pdf'
```

## Available Tools

The MCP server provides the following tools for interacting with the Lynx Reservations system:

### MCP Tools

#### 1. `file_search_by_party_name`
**Description:** Retrieve file from party name  
**Usage:** Search for files using a customer's last name

#### 2. `file_search_by_file_reference` 
**Description:** Retrieve file from file reference  
**Usage:** Search for files using the Lynx file reference (e.g., FTXXXXXXXXX)

#### 3. `retrieve_itinerary`
**Description:** Retrieve file itinerary  
**Usage:** Get detailed itinerary information for a specific file

#### 4. `retrieve_file_documents`
**Description:** Retrieve file documents from transaction reference  
**Usage:** Get documents associated with a specific transaction

#### 5. `attachment_upload`
**Description:** Upload attachment for using with file document  
**Usage:** Upload binary files (PDFs, images, etc.) to be associated with documents

> **Note:** This tool currently doesn't scale well due to the attachment increasing the content window too much and hitting OpenAI rate limits, please use the REST endpoint meanwhile for your orchestration!

#### 6. `file_document_save`
**Description:** Save file document details  
**Usage:** Create or update documents at the file level

#### 7. `transaction_document_save`
**Description:** Save transaction document details  
**Usage:** Create or update documents at the transaction level

### REST Endpoints

#### POST `/attachmentUpload`
**Description:** REST endpoint for uploading file attachments  
**Authentication:** Requires Bearer token (same as MCP server)  
**Content-Type:** `multipart/form-data`  
**Parameters:**
- `file` (required): The file to upload (max 32MB)
- `fileId` (required): The Lynx file identifier

**Response:** JSON with the attachment URL
```json
{
  "attachmentUrl": "/documents/file/f16476987/d20250708231038.pdf"
}
```

**Example Usage:**
```bash
curl -X POST http://localhost:9600/attachmentUpload \
  -H "Authorization: Bearer YOUR_BEARER_TOKEN" \
  -F "file=@document.pdf" \
  -F "fileId=12345"
```

**Error Responses:**
- `400 Bad Request`: Missing required parameters or invalid file
- `401 Unauthorized`: Invalid or missing Bearer token
- `500 Internal Server Error`: Server-side processing error