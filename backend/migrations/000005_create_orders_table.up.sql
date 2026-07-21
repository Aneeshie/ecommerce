CREATE TYPE order_status AS ENUM ('pending', 'paid', 'cancelled');

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    status order_status NOT NULL DEFAULT 'pending',
    total_price BIGINT NOT NULL CHECK (total_price >= 0),

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_orders_user_id ON orders(user_id);