// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"
import { OpenAIEmbeddings } from "npm:@langchain/openai";
import { SupabaseVectorStore } from "npm:@langchain/community/vectorstores/supabase";
import { createClient } from 'jsr:@supabase/supabase-js@2.49.5';

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

Deno.serve(async (req) => {
  // Grab the user's query from the JSON payload
  const { query, filterByBookingRef = null } = await req.json();
  if (!query) {
    throw new Error("query is required");
  }

  const client = createClient(supabaseUrl, supabaseServiceRoleKey);

  const embeddings = new OpenAIEmbeddings({
    apiKey: openaiApiKey,
    model: 'text-embedding-3-small',
    dimensions: 1536
  });

  const vectorStore = new SupabaseVectorStore(embeddings, {
    client,
    tableName: "emails",
    queryName: "semantic_search",
  });

  const funcFilterOnBookingRef = (rpc) => {
    if (!!filterByBookingRef) {
      rpc.filter("metadata->>bookingReference::string", "ilike", `%${filterByBookingRef}%`);
    }
    return rpc;
  }

  const results = await vectorStore.similaritySearchWithScore(
    query,
    2,
    funcFilterOnBookingRef
  );

  return new Response(JSON.stringify(results.map(([{ pageContent: emailContent, metadata }, score]) => ({ emailContent, metadata, score }))), {
    headers: {
      'Content-Type': 'application/json'
    }
  });
})

/* To invoke locally:

  1. Run `supabase start` (see: https://supabase.com/docs/reference/cli/supabase-start)
  2. Make an HTTP request:

  curl -i --location --request POST 'http://127.0.0.1:54321/functions/v1/semantic_search_function' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0' \
    --header 'Content-Type: application/json' \
    --data '{"name":"Functions"}'

*/
