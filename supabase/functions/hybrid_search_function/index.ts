// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"
import { OpenAIEmbeddings } from "npm:@langchain/openai";
import { SupabaseVectorStore } from "npm:@langchain/community/vectorstores/supabase";
import { createClient, SupabaseClient } from 'jsr:@supabase/supabase-js';
import { CohereClient } from "npm:cohere-ai";

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

const cohereApiKey = Deno.env.get('COHERE_API_KEY');
if (!cohereApiKey) {
  throw new Error("COHERE_API_KEY environment variable is required");
}

const EMBEDDING_MODEL = "text-embedding-3-small";
const EMBEDDING_SIZE = 1536;

const TABLE_NAME = "emails";
const SIMILARITY_QUERY_NAME = "semantic_search";
const KEYWORD_QUERY_NAME = "kw_hybrid_search";
const SIMILARITY_K = 4;
const KEYWORD_K = 3;

interface SearchResult {
  pageContent: string;
  metadata: any;
  score?: number;
}

interface RerankedResult {
  index: number;
  relevanceScore: number;
}

async function rerankWithCohere(
  query: string,
  documents: SearchResult[]
): Promise<SearchResult[]> {
  if (documents.length === 0) {
    return documents;
  }

  const cohere = new CohereClient({
    token: cohereApiKey,
  });

  try {
    const rerankResponse = await cohere.rerank({
      query: query,
      documents: documents.map(doc => doc.pageContent),
      topN: documents.length,
      model: "rerank-english-v3.0",
    });

    // Create a map of original index to reranked result
    const rerankedMap = new Map<number, RerankedResult>();
    rerankResponse.results.forEach((result: any) => {
      rerankedMap.set(result.index, {
        index: result.index,
        relevanceScore: result.relevanceScore
      });
    });

    // Reorder documents based on Cohere's reranking
    const rerankedDocuments = documents.map((doc, originalIndex) => {
      const reranked = rerankedMap.get(originalIndex);
      return {
        ...doc,
        score: reranked ? reranked.relevanceScore : doc.score || 0
      };
    });

    // Sort by Cohere's relevance score
    return rerankedDocuments.sort((a, b) => (b.score || 0) - (a.score || 0));
  } catch (error) {
    console.error("Cohere reranking failed:", error);
    // Fallback to original ranking if Cohere fails
    return documents.sort((a, b) => (b.score || 0) - (a.score || 0));
  }
}

async function customHybridSearch(
  client: SupabaseClient,
  vectorStore: SupabaseVectorStore,
  query: string,
  filterByFileReference: string | null
): Promise<SearchResult[]> {
  const filter = filterByFileReference ? {
    fileReference: filterByFileReference
  } : {};

  const funcFilterOnBookingRef = (rpc: any) => {
    if (filterByFileReference) {
      rpc.filter("metadata->>fileReference::string", "ilike", `%${filterByFileReference}%`);
    }
    return rpc;
  }

  // Call both similarity and keyword search functions
  const [similarityResults, keywordResults] = await Promise.all([
    await vectorStore.similaritySearchWithScore(
      query,
      SIMILARITY_K,
      funcFilterOnBookingRef
    ),
    client.rpc(KEYWORD_QUERY_NAME, {
      query_text: query,
      match_count: KEYWORD_K,
      filter: filter,
    })
  ]);

  // Combine and deduplicate results
  const allResults = new Map<string, SearchResult>();

  // Process similarity results
  if (Array.isArray(similarityResults)) {
    console.log("similarityResults length", similarityResults.length);
    similarityResults.forEach(([{ pageContent, metadata }, score]: any) => {
      const key = `${pageContent}-${JSON.stringify(metadata)}`;
      allResults.set(key, {
        pageContent,
        metadata,
        score
      });
    });
  }

  // Process keyword results
  if (keywordResults?.data && Array.isArray(keywordResults.data)) {
    console.log("keywordResults length", keywordResults.data.length);
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

  const combinedResults = Array.from(allResults.values());

  // Apply Cohere reranking to improve the final ranking
  return await rerankWithCohere(query, combinedResults);
}

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
    dimensions: EMBEDDING_SIZE,
  });

  const vectorStore = new SupabaseVectorStore(embeddings, {
    client,
    tableName: TABLE_NAME,
    queryName: SIMILARITY_QUERY_NAME,
  });

  const results = await customHybridSearch(client, vectorStore, query, filterByFileReference);

  console.log("results length:", results.length);

  return new Response(
    JSON.stringify(
      results
        .filter(({ score }) => !!score && score > 0)
        .map(({ pageContent: emailContent, metadata, score }) => ({ emailContent, metadata, score }))
    ), {
    headers: {
      'Content-Type': 'application/json'
    }
  });
})
