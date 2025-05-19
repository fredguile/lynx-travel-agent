import {
  SupabaseFilterRPCCall,
  SupabaseVectorStore,
} from "@langchain/community/vectorstores/supabase";
import { ChatOpenAI, OpenAIEmbeddings } from "@langchain/openai";
import { createClient } from "@supabase/supabase-js";

const privateKey = process.env.SUPABASE_SERVICE_ROLE_KEY;
if (!privateKey) {
  throw new Error("SUPABASE_SERVICE_ROLE_KEY environment variable is required");
}
const url = "https://tzdazzxqjilbkrxufnbu.supabase.co";

export const run = async () => {
  const client = createClient(url, privateKey);

  const query = "ground transportation";

  const embeddings = new OpenAIEmbeddings({
    model: "text-embedding-3-small",
  }); // requires process.env.OPENAI_API_KEY

  const queryEmbedding = await embeddings.embedQuery(query);

  // const funcFilterOnBookingRef: SupabaseFilterRPCCall = (rpc) =>
  //   rpc.filter("metadata->>bookingReference::string", "eq", "FTCA5250128");

  const { data: emails } = await client.rpc('hybrid_search', { query_text: query, query_embedding: queryEmbedding, match_count: 10 })

  console.log(emails);
};

run();
