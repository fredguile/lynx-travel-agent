// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"
import { OpenAIEmbeddings } from "npm:@langchain/openai";
import { SupabaseVectorStore } from "npm:@langchain/community/vectorstores/supabase";
import { createClient } from 'jsr:@supabase/supabase-js@2.50.0';

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

const TABLE_NAME = "emails";
const QUERY_NAME = "semantic_search";

Deno.serve(async (req) => {
    // Grab the user's query from the JSON payload
    const { query, filterByBookingRef = null, filterByBookingConfirmationId = null } = await req.json();
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
        tableName: TABLE_NAME,
        queryName: QUERY_NAME,
    });

    const funcFilterOnBookingRef = (rpc: any) => {
        if (filterByBookingRef) {
            rpc.filter("metadata->>bookingReference::string", "ilike", `%${filterByBookingRef}%`);
        }
        if (filterByBookingConfirmationId) {
            rpc.filter("metadata->>bookingConfirmationId::string", "like", `%${filterByBookingConfirmationId}%`);
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