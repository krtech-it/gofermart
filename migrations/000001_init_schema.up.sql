CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    number      TEXT PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users (id),
    status      order_status NOT NULL DEFAULT 'NEW',
    accrual     NUMERIC(18, 2),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX orders_user_id_idx ON orders (user_id);

CREATE TABLE withdrawals (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID        NOT NULL REFERENCES users (id),
    order_number TEXT        NOT NULL,
    sum          NUMERIC(18, 2) NOT NULL CHECK (sum > 0),
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX withdrawals_user_id_idx ON withdrawals (user_id);
