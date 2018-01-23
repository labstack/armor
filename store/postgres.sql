create table if not exists plugins (
  name text primary key,
  config json not null,
  enabled boolean not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create table if not exists proxy_targets (
  name text not null,
  url text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
)
