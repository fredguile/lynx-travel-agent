// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"
import { OpenAIEmbeddings } from "npm:@langchain/openai";
import { SupabaseVectorStore } from "npm:@langchain/community/vectorstores/supabase";
import { createClient } from 'jsr:@supabase/supabase-js';

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

const TABLE_NAME = "emails";
const QUERY_NAME = "semantic_search";
const NUM_RESULTS = 4;

Deno.serve(async (req) => {
    // Grab the user's query from the JSON payload
    const { query, filterByFileReference = null } = await req.json();
    if (!query) {
        throw new Error("query is required");
    }

    console.log("query:", query);
    console.log("filterByFileReference:", filterByFileReference);

    const client = createClient(supabaseUrl, supabaseServiceRoleKey);

    const embeddings = new OpenAIEmbeddings({
        apiKey: openaiApiKey,
        model: EMBEDDING_MODEL,
        dimensions: EMBEDDING_SIZE
    });

    const vectorStore = new SupabaseVectorStore(embeddings, {
        client,
        tableName: TABLE_NAME,
        queryName: QUERY_NAME,
    });

    const funcFilterOnBookingRef = (rpc: any) => {
        if (filterByFileReference) {
            rpc.filter("metadata->>fileReference::string", "ilike", `%${filterByFileReference}%`);
        }
        return rpc;
    }

    const results = await vectorStore.similaritySearchWithScore(
        query,
        NUM_RESULTS,
        funcFilterOnBookingRef
    );

    console.log("results length:", results.length);

    return new Response(
        JSON.stringify(
            results
                .filter(([_, score]) => score > 0)
                .map(([{ pageContent: emailContent, metadata }, score]) => ({ emailContent, metadata, score })
                )
        ), {
        headers: {
            'Content-Type': 'application/json'
        }
    });
})
