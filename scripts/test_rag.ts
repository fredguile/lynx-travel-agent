import {
  SupabaseFilterRPCCall,
  SupabaseVectorStore,
} from "@langchain/community/vectorstores/supabase";
import { ChatOpenAI, OpenAIEmbeddings } from "@langchain/openai";
import { createClient } from "@supabase/supabase-js";

const privateKey =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InR6ZGF6enhxamlsYmtyeHVmbmJ1Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc0NzU5NzYyNywiZXhwIjoyMDYzMTczNjI3fQ.ML1e53mqeVebS9BfHntbxDzYPk02NTw24P4LO7yA7WE";
const url = "https://tzdazzxqjilbkrxufnbu.supabase.co";

export const run = async () => {
  const client = createClient(url, privateKey);

  const embeddings = new OpenAIEmbeddings({
    model: "text-embedding-3-small",
  });

  const store = new SupabaseVectorStore(embeddings, {
    client,
    tableName: "emails",
    queryName: "match_emails",
  });

  const funcFilterOnBookingRef: SupabaseFilterRPCCall = (rpc) =>
    rpc.filter("metadata->>booking_reference::string", "eq", "1708634");

  const result = await store.similaritySearchWithScore(
    "mercure",
    4,
    funcFilterOnBookingRef,
  );

  // Format the context from the most relevant results
  const context = result
    .map(
      ([doc, score], i) =>
        `Result ${i + 1} (score: ${score}):\n${doc.pageContent}`,
    )
    .join("\n\n");

  // The user prompt
  const userPrompt = "Should I book a Mercure hotel for these customers?";

  // Compose the full prompt for the LLM
  const fullPrompt = `Context:\n${context}\n\nQuestion: ${userPrompt}`;

  // Import ChatOpenAI
  const chat = new ChatOpenAI({
    modelName: "gpt-3.5-turbo",
    temperature: 0,
  });

  // Send the prompt to the LLM
  const response = await chat.invoke([
    [
      "system",
      "You are a helpful assistant. Answer the user's question using only the provided context. If the context is insufficient, say so.",
    ],
    ["human", fullPrompt],
  ]);

  console.log("LLM answer:\n", response.content);
};

run();
