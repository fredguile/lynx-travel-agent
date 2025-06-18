# panpac-helper-browser-extension

Firefox/Chrome extension to use as helper for Pan PAC day-to-day tasks.

> **Notice:** This extension is currently only tested with **Firefox**.

## Table of Content

- [browser-extension](./browser-extension/)
  - **A Firefox/Chrome extension** to assist with Pan PAC day-to-day tasks. Includes build and development instructions for running the extension locally.
- [docker-compose](./docker-compose/)
  - **Docker Compose configuration** for running n8n.
- [lynx-mcp-server](./lynx-mcp-server/)
  - **MCP Server for Lynx Reservations written in golang**
- [n8n_templates](./n8n_templates/)
  - **Automation workflow templates** for n8n, including pipelines and chatbot integrations to streamline Pan PAC processes.
- [supabase](./supabase/)
  - **Backend and database resources** for Pan PAC, featuring serverless functions (for hybrid and semantic search), and SQL schemas (with and without full-text search) to support advanced data operations.

##  Supabase Commands

- `yarn login`: Login to Supabase
- `yarn start`: Start supabase local environment
- `yarn stop`: Stop supabase local environment
- `yarn status`: Check status of supabase local environment
- `yarn db-link`: Link local database to remote project
- `yarn db-reset`: Reset local database
- `yarn db-pull`: Pull remote database schema as migration (see
  supabase/migrations folder)
- `yarn db-migration`: Manage database migrations
- `yarn db-dump`: Dump from the production database (example: with `--data-only`)
- `yarn edge-functions`: Manage Supabase edge functions
- `yarn secrets`: Manage Supabase secrets