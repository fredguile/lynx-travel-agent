create or replace function hybrid_search(
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

-- Create a function to keyword search for documents
-- example call: select * from kw_hybrid_search('flight details', 4, jsonb_build_object('fileReference', 'FT1740618'));
create or replace function kw_hybrid_search(
  query_text text,
  match_count int,
  filter jsonb DEFAULT '{}'::jsonb
)
returns table (
  id bigint, 
  content text, 
  metadata jsonb, 
  similarity real
)
language plpgsql 
as $$
#variable_conflict use_column
begin
  return query 
  execute format('select id, content, metadata, ts_rank(to_tsvector(content), phraseto_tsquery($1)) as similarity
    from emails
    where metadata @> $3
    and to_tsvector(content) @@ phraseto_tsquery($1)
    order by similarity desc
    limit $2'
  ) using query_text, match_count, filter;
end;
$$;
