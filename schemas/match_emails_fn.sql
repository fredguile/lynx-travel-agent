CREATE OR REPLACE FUNCTION public.match_emails(query_embedding vector, match_count integer DEFAULT NULL::integer, filter jsonb DEFAULT '{}'::jsonb)
 RETURNS TABLE(id bigint, content text, metadata jsonb, similarity double precision)
 LANGUAGE plpgsql
AS $function$#variable_conflict use_column

begin  

return query  select    id,    content,    metadata,    1 - (emails.embedding <=> query_embedding) as similarity  from emails where metadata @> filter  order by emails.embedding <=> query_embedding  limit match_count;

end;$function$