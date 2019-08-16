create table orders
(
    id          varchar not null
        constraint orders_pk
            primary key,
    variant_sku varchar not null
        constraint orders_inventory_variants_sku_fk
            references inventory_variants
            on update restrict on delete restrict,
    quantity    int       default 0 not null,
    price       float     default 0,
    receipt     varchar   default null,
    created_at  timestamp default current_timestamp not null,
    updated_at  timestamp default current_timestamp not null,
    deleted_at  timestamp default null
);

create table transactions
(
    id          varchar not null
        constraint transactions_pk
            primary key,
    variant_sku varchar not null
        constraint transactions_inventory_variants_sku_fk
            references inventory_variants
            on update restrict on delete restrict,
    order_id    varchar   default null
        constraint transactions_orders_id_fk
            references orders
            on update restrict on delete restrict,
    type        varchar not null,
    quantity    int       default 0 not null,
    price       float     default 0,
    created_at  timestamp default current_timestamp not null,
    updated_at  timestamp default current_timestamp not null,
    deleted_at  timestamp default null
);
