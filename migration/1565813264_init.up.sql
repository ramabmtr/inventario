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
    sku          varchar not null
        constraint inventory_variant_pk
            primary key,
    inventory_id varchar not null
        constraint inventory_variants_inventories_id_fk
            references inventories
            on update cascade on delete cascade,
    name         varchar   default null,
    size         varchar   default null,
    color        varchar   default null,
    quantity     integer   default 0 not null,
    created_at   timestamp default current_timestamp not null,
    updated_at   timestamp default current_timestamp not null,
    deleted_at   timestamp default null
);
