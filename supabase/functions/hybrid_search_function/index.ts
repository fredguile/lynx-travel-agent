// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"
import { OpenAIEmbeddings } from "npm:@langchain/openai";
import { createClient, SupabaseClient } from 'jsr:@supabase/supabase-js@2.50.2';

const supabaseUrl = Deno.env.get('SUPABASE_URL');
if (!supabaseUrl) {
  throw new Error("SUPABASE_URL environment variable is required");
}

const supabaseServiceRoleKey = Deno.env.get('SUPABASE_SERVICE_ROLE_KEY');
if (!supabaseServiceRoleKey) {
  throw new Error("SUPABASE_SERVICE_ROLE_KEY environment variable is required");
}

const openaiApiKey = Deno.env.get('OPENAI_API_KEY');
if (!openaiApiKey) {
  throw new Error("OPENAI_API_KEY environment variable is required");
}

const EMBEDDING_MODEL = "text-embedding-3-small";
const EMBEDDING_SIZE = 1536;

const SIMILARITY_QUERY_NAME = "hybrid_search";
const KEYWORD_QUERY_NAME = "kw_hybrid_search";
const SIMILARITY_K = 3;
const KEYWORD_K = 3;

interface SearchResult {
  pageContent: string;
  metadata: any;
  score?: number;
}

async function customHybridSearch(
  client: SupabaseClient,
  embeddings: OpenAIEmbeddings,
  query: string,
  filterByFileReference: string | null
): Promise<SearchResult[]> {
  // Generate embedding for the query
  const queryEmbedding = await embeddings.embedQuery(query);

  const filter = filterByFileReference ? {
    fileReference: filterByFileReference
  } : {};

  // Call both similarity and keyword search functions
  const [similarityResults, keywordResults] = await Promise.all([
    client.rpc(SIMILARITY_QUERY_NAME, {
      query_embedding: queryEmbedding,
      match_count: SIMILARITY_K,
      filter: filter
    }),
    client.rpc(KEYWORD_QUERY_NAME, {
      query_text: query,
      match_count: KEYWORD_K,
      filter: filter,
    })
  ]);

  // Combine and deduplicate results
  const allResults = new Map<string, SearchResult>();

  // Process similarity results
  if (similarityResults?.data && Array.isArray(similarityResults.data)) {
    similarityResults.data.forEach((result: any) => {
      const key = `${result.content}-${JSON.stringify(result.metadata)}`;
      allResults.set(key, {
        pageContent: result.content,
        metadata: result.metadata,
        score: result.similarity
      });
    });
  }
  
  // Process keyword results
  if (keywordResults?.data && Array.isArray(keywordResults.data)) {
    keywordResults.data.forEach((result: any) => {
      const key = `${result.content}-${JSON.stringify(result.metadata)}`;
      if (!allResults.has(key)) {
        allResults.set(key, {
          pageContent: result.content,
          metadata: result.metadata,
          score: result.similarity
        });
      }
    });
  }

  return Array.from(allResults.values()).sort((a, b) => (b.score || 0) - (a.score || 0));
}

Deno.serve(async (req) => {
  // Grab the user's query from the JSON payload
  const { query, filterByFileReference = null } = await req.json();
  if (!query) {
    throw new Error("query is required");
  }

  const client = createClient(supabaseUrl, supabaseServiceRoleKey);

  const embeddings = new OpenAIEmbeddings({
    apiKey: openaiApiKey,
    model: EMBEDDING_MODEL,
    dimensions: EMBEDDING_SIZE
  });

  const results = await customHybridSearch(client, embeddings, query, filterByFileReference);

  return new Response(
    JSON.stringify(
      results
        .map(({ pageContent: emailContent, metadata, score }) => ({ emailContent, metadata, score }))
    ), {
    headers: {
      'Content-Type': 'application/json'
    }
  });
})
