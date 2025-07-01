# Hybrid Search Function with Cohere Reranking

This Supabase Edge Function implements a hybrid search system that combines vector similarity search and keyword search, then uses Cohere's reranking API to improve result relevance.

## Features

- **Vector Similarity Search**: Uses OpenAI embeddings to find semantically similar content
- **Keyword Search**: Uses PostgreSQL full-text search for exact keyword matching
- **Cohere Reranking**: Applies Cohere's rerank-english-v3.0 model to improve final ranking
- **Deduplication**: Removes duplicate results from both search methods
- **Filtering**: Supports filtering by file reference

## Environment Variables

The following environment variables are required:

- `SUPABASE_URL`: Your Supabase project URL
- `SUPABASE_SERVICE_ROLE_KEY`: Your Supabase service role key
- `OPENAI_API_KEY`: Your OpenAI API key for embeddings
- `COHERE_API_KEY`: Your Cohere API key for reranking

## API Usage

### Request

```json
{
  "query": "your search query",
  "filterByFileReference": "optional_file_reference"
}
```

### Response

```json
[
  {
    "emailContent": "matched content",
    "metadata": {
      "fileReference": "file_ref",
      "otherMetadata": "value"
    },
    "score": 0.95
  }
]
```

## How It Works

1. **Query Processing**: The user's query is embedded using OpenAI's text-embedding-3-small model
2. **Parallel Search**: Both similarity and keyword searches are executed simultaneously
3. **Result Combination**: Results from both searches are combined and deduplicated
4. **Cohere Reranking**: The combined results are reranked using Cohere's rerank-english-v3.0 model
5. **Final Ranking**: Results are sorted by Cohere's relevance score

## Error Handling

- If Cohere reranking fails, the function falls back to the original similarity-based ranking
- All API errors are logged for debugging purposes

## Performance Considerations

- The function uses `Promise.all()` for parallel execution of similarity and keyword searches
- Cohere reranking is applied to the combined results to minimize API calls
- Results are limited by `SIMILARITY_K` and `KEYWORD_K` constants (both set to 3 by default)

## Dependencies

- `@langchain/openai`: For OpenAI embeddings
- `@supabase/supabase-js`: For Supabase client
- `cohere-ai`: For Cohere reranking API 