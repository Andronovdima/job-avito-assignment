CREATE TABLE IF NOT EXISTS balances (
    id bigserial not null primary key,
    user_id bigint unique not null,
    balance bigint DEFAULT 0
);