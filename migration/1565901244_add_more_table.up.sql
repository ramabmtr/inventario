create table orders
(
    id          varchar not null
        constraint orders_pk
            primary key,
    variant_sku varchar not null
        constraint orders_inventory_variants_sku_fk
            references inventory_variants
            on update restrict on delete restrict,
    quantity    int     not null,
    price       float   not null,
    receipt     varchar   default null,
    created_at  timestamp default current_timestamp not null,
    updated_at  timestamp default current_timestamp not null,
    deleted_at  timestamp default null
);

create table order_transactions
(
    id         varchar not null
        constraint transaction_items_pk
            primary key,
    order_id   varchar not null
        constraint order_transactions_orders_id_fk
            references orders
            on update restrict on delete restrict,
    quantity   int     not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null,
    deleted_at timestamp default null
);

create table transactions
(
    id         varchar not null
        constraint transactions_pk
            primary key,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null,
    deleted_at timestamp default null
);

create table transaction_items
(
    id             varchar not null
        constraint transaction_items_pk
            primary key,
    transaction_id varchar not null
        constraint transaction_items_transactions_id_fk
            references transactions
            on update restrict on delete restrict,
    variant_sku    varchar not null
        constraint transaction_items_inventory_variants_sku_fk
            references inventory_variants
            on update restrict on delete restrict,
    quantity       int     not null,
    price          float   not null,
    created_at     timestamp default current_timestamp not null,
    updated_at     timestamp default current_timestamp not null,
    deleted_at     timestamp default null
);
