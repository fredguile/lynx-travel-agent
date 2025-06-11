create or replace function hybrid_search(
  query_text text,
  query_embedding vector(1536),
  match_count int,
  filter jsonb DEFAULT '{}'::jsonb,
  full_text_weight float = 1,
  semantic_weight float = 1,
  rrf_k int = 50
)
returns table (
  id bigint,
  created_at timestamp,
  content text,
  metadata jsonb,
  full_text_weight float,
  semantic_weight float
)
language sql
as $$
with log_entry as (
  insert into function_logs(function_name, message)
  values ('hybrid_search', 'filter=' || filter::text)
  returning 1
),
full_text as (
  select
    id,
    -- Note: ts_rank_cd is not indexable but will only rank matches of the where clause
    -- which shouldn't be too big
    row_number() over(order by ts_rank_cd(fts, websearch_to_tsquery(query_text)) desc) as rank_ix
  from
    emails
  where
    fts @@ websearch_to_tsquery(query_text)
    and (metadata @> filter OR filter = '{}'::jsonb)
  order by rank_ix
  limit least(match_count, 30) * 2
),
semantic as (
  select
    id,
    row_number() over (order by embedding <#> query_embedding) as rank_ix
  from
    emails
  where
    (metadata @> filter OR filter = '{}'::jsonb)
  order by rank_ix
  limit least(match_count, 30) * 2
)
select
  emails.id,
  emails.created_at,
  emails.content,
  emails.metadata,
  full_text_weight,
  semantic_weight
from
  full_text
  full outer join semantic
    on full_text.id = semantic.id
  join emails
    on coalesce(full_text.id, semantic.id) = emails.id
order by
  coalesce(1.0 / (rrf_k + full_text.rank_ix), 0.0) * full_text_weight +
  coalesce(1.0 / (rrf_k + semantic.rank_ix), 0.0) * semantic_weight
  desc
limit
  least(match_count, 30)
$$;