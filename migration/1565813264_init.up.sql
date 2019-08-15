create table inventories
(
    id         varchar not null
        constraint inventory_pk
            primary key,
    name       varchar not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null,
    deleted_at timestamp default null
);

create table inventory_variants
(
    id         varchar not null
        constraint inventory_variant_pk
            primary key,
    sku        varchar   default null,
    name       varchar   default null,
    size       varchar   default null,
    color      varchar   default null,
    quantity   integer   default 0 not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null,
    deleted_at timestamp default null
);