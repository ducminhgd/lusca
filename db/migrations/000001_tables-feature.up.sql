create table if not exists feature (
    id UUID primary key,
    name TEXT not null,
    description TEXT,
    status INT2 not null default 1
);

create table if not exists strategy_kv (
    id UUID primary key,
    feature_id UUID not null,
    environment JSONB,
    type INT2 not null default 1,
    key TEXT not null,
    value JSONB
);

create index idx_strategy_key_key on strategy_kv (key);