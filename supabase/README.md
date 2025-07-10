# supabase

Supabase assets and configuration for the Lynx Travel Agent.

## Prerequisites

The `vector` extension must be enabled.

## Folder Structure

- `functions/`: Contains Supabase edge functions for search capabilities.
  - `semantic_search_function/`: TypeScript/Deno code for semantic search.
  - `hybrid_search_function/`: TypeScript/Deno code for hybrid search.

- `migrations/`: Contains database migration files that track schema changes over time.
  These files are automatically generated and managed by Supabase CLI to ensure
  database schema version control and consistent deployments across environments.

- `user_schemas/`: Contains SQL schema files for database setup.
  - `with_fts/`: Schemas with full-text search (FTS) enabled (e.g.,
    `emails.sql`, `hybrid_search.sql`).
  - `without_fts/`: Schemas without FTS (e.g., `semantic_search.sql`,
    `emails_no_fts.sql`).
