  drop function if exists semantic_search_dev(vector,integer,jsonb);

  create or replace function semantic_search_dev(
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
      1 - (emails_dev.embedding <=> query_embedding) as similarity
    from emails_dev
    where metadata @> filter
    order by emails_dev.embedding <=> query_embedding
    limit match_count;
  end;
  $$;

  drop function if exists semantic_search_prod(vector,integer,jsonb);

  create or replace function semantic_search_prod(
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
      1 - (emails_prod.embedding <=> query_embedding) as similarity
    from emails_prod
    where metadata @> filter
    order by emails_prod.embedding <=> query_embedding
    limit match_count;
  end;
  $$;
