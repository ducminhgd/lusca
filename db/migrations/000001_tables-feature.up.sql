create table if not exists "feature" (
    id UUID primary key,
    name TEXT not null,
    description TEXT,
    status INT2 not null default 1,
    created_at TIMESTAMPTZ not null default now(),
    updated_at TIMESTAMPTZ
);
create unique index unq_feature_name on "feature" (name);

create table if not exists "strategy" (
    id UUID primary key,
    feature_id UUID not null,
    environment JSONB,
    data JSONB
);

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
create unique index unq_collection_detail_collection_id_value on "collection_detail" (collection_id, value);

create table if not exists "feature_collection" (
    feature_id UUID not null,
    collection_id UUID not null,
    PRIMARY KEY (feature_id, collection_id)
);


-- Test Data
-- INSERT INTO public.collection (id,"name",description) VALUES ('8a735d82-80cf-4b62-87f8-673a995771bf'::uuid,'test','test collection');
-- INSERT INTO public.collection_detail (collection_id,value) VALUES ('8a735d82-80cf-4b62-87f8-673a995771bf'::uuid,'user_id_1'), ('8a735d82-80cf-4b62-87f8-673a995771bf'::uuid,'user_id_2');

