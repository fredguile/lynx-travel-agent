set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.hybrid_search(query_embedding vector, match_count integer DEFAULT NULL::integer, filter jsonb DEFAULT '{}'::jsonb)
 RETURNS TABLE(id bigint, content text, metadata jsonb, similarity double precision)
 LANGUAGE plpgsql
AS $function$
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
$function$
;

CREATE OR REPLACE FUNCTION public.kw_hybrid_search(query_text text, match_count integer, filter jsonb DEFAULT '{}'::jsonb)
 RETURNS TABLE(id bigint, content text, metadata jsonb, similarity real)
 LANGUAGE plpgsql
AS $function$
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
$function$
;


