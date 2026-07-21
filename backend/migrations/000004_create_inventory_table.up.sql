
CREATE TABLE inventories(
   product_id UUID PRIMARY KEY REFERENCES products(id) ON DELETE RESTRICT,

   quantity INTEGER NOT NULL CHECK (quantity >= 0),

   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL
);