# lynx-travel-agent

A pet project building the best AI agent to work with lynx-reservations. Work in progress. Usage of RAG knowledge paired with helpful Firefox/Chrome browser extension.

> **Notice:** Extension is currently only tested with **Firefox**.

## Terms of use

[![License: CC BY-NC-SA 4.0](https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

This project is licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License. See [LICENSE](LICENSE) for details.

## Table of Content

- [browser-extension](./browser-extension/)
  - **Firefox/Chrome extension** to use as Lynx Travel Agent. Includes build and development instructions for running the extension locally.
- [docker-compose](./docker-compose/)
  - **Docker Compose configuration** for running n8n.
- [lynx-mcp-server](./lynx-mcp-server/)
  - **MCP Server for Lynx Reservations written in golang**
- [n8n_templates](./n8n_templates/)
  - **Automation workflow templates** for n8n, including RAG pipeline and chatbot integration.
- [supabase](./supabase/)
  - **Backend and database resources** on Supabase, featuring serverless functions (for hybrid and semantic search) and SQL schemas (with and without full-text search) to support advanced data operations.

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
