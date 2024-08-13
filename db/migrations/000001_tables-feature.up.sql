create table if not exists "feature" (
    id UUID primary key,
    name TEXT not null,
    description TEXT,
    status INT2 not null default 1,
    created_at TIMESTAMPTZ not null default now(),
    updated_at TIMESTAMPTZ
);

create table if not exists "strategy_kv" (
    id UUID primary key,
    feature_id UUID not null,
    environment JSONB,
    type INT2 not null default 1,
    key TEXT not null,
    value JSONB
);
create index idx_strategy_key_key on "strategy_kv" (key);

create table if not exists "collection" (
    id UUID primary key,
    name TEXT not null,
    description TEXT
);

create table if not exists "collection_detail" (
    collection_id UUID not null,
    value TEXT NOT NULL
);
create index idx_collection_detail_value on "collection_detail" (value);

create table if not exists "feature_collection" (
    feature_id UUID not null,
    collection_id UUID not null,
    PRIMARY KEY (feature_id, collection_id)
);