
CREATE TYPE product_status AS ENUM ('ACTIVE', 'ARCHIVED');

CREATE TABLE products (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  price BIGINT NOT NULL CHECK (price >= 0),
  status product_status NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
