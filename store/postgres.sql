create table if not exists plugins (
  id text primary key,
  name text not null,
  host text not null,
  path text not null,
  config jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  unique (name, host, path)
);

-- create table if not exists proxy_targets (
--   name text not null,
--   url text not null,
--   created_at timestamptz not null,
--   updated_at timestamptz not null
-- )
