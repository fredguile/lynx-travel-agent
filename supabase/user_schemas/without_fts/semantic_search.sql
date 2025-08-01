create or replace function semantic_search(
    query_embedding vector(1536),
    match_count int DEFAULT null,
    filter jsonb DEFAULT '{}'::jsonb
  ) 
  returns table (
    id bigint,
    content text,
    metadata jsonb,
    similarity float
  )
  language plpgsql
  as $$
  #variable_conflict use_column
  begin
    return query
    select
      id,
      content,
      metadata,
      1 - (emails.embedding <=> query_embedding) as similarity
    from emails
    where metadata @> filter
    order by emails.embedding <=> query_embedding
    limit match_count;
  end;
  $$;
  
