# panpac-supabase

Supabase assets and configuration for the PAN Pac extension.

## Prerequisites

The `vector` extension must be enabled.

## Folder Structure

- `functions/`: Contains Supabase edge functions for search capabilities.
  - `semantic_search_function/`: TypeScript/Deno code for semantic search.
  - `hybrid_search_function/`: TypeScript/Deno code for hybrid search.
  
- `schemas/`: Contains SQL schema files for database setup.
  - `with_fts/`: Schemas with full-text search (FTS) enabled (e.g., `emails.sql`, `hybrid_search.sql`).
  - `without_fts/`: Schemas without FTS (e.g., `semantic_search.sql`, `emails_no_fts.sql`).

