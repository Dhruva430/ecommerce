CREATE TABLE
  orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user" (id),
    address JSONB NOT NULL,
    seller_id BIGINT REFERENCES seller (id),
    address_id BIGINT REFERENCES address (id),
    total_amount FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    status order_status NOT NULL DEFAULT 'PENDING',
    payment_status payment_status NOT NULL DEFAULT 'PENDING'
  );

CREATE TABLE
  order_product (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES product (id),
    seller_id BIGINT NOT NULL REFERENCES seller (id),
    amount INT NOT NULL,
    variant_id BIGINT NOT NULL REFERENCES product_variant (id),
    order_id BIGINT REFERENCES orders (id)
  );

CREATE TABLE
  payment (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE REFERENCES orders (id),
    transaction_id TEXT UNIQUE,
    payment_gateway TEXT,
    amount FLOAT NOT NULL,
    status payment_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW ()
  );