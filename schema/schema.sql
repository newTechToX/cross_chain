create table "common_cross_chain" (
    "id" bigserial not null,
    "match_id" int8,
    "chain" varchar not null,
    "number" int8 not null,
    "ts" timestamp not null,
    "index" int8 not null,
    "hash" varchar not null,
    "action_id" int8 not null,
    "project" varchar not null,
    "contract" varchar not null,
    "direction" varchar not null,
    "from_chain_id" numeric(256),
    "from_address" varchar not null,
    "to_chain_id" numeric(256),
    "to_address" varchar not null,
    "token" varchar not null,
    "amount" numeric(256),
    "match_tag" varchar not null,
    "detail" json,
    primary key (
        "project",
        "chain",
        "number",
        "hash",
        "index",
        "action_id"
    )
);

create index on common_cross_chain (id);

create index on common_cross_chain (match_id);

create index on common_cross_chain (hash);

create index on common_cross_chain (match_tag);

create index on common_cross_chain (chain, project, match_tag);

create index on common_cross_chain (chain, project);

create index on common_cross_chain (from_address);

create index on common_cross_chain (to_address);

create index on common_cross_chain (chain, from_address);

create index on common_cross_chain (chain, to_address);

create index on common_cross_chain (ts desc);