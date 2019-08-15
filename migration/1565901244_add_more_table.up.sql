create table transactions
(
    id          varchar not null
        constraint transactions_pk
            primary key,
    variant_sku varchar not null
        constraint transactions_inventory_variants_sku_fk
            references inventory_variants
            on update restrict on delete restrict,
    type        varchar not null,
    quantity    int       default 0 not null,
    price       float     default 0,
    created_at  timestamp default current_timestamp not null,
    updated_at  timestamp default current_timestamp not null,
    deleted_at  timestamp default null
);
