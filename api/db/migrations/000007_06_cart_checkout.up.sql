CREATE TABLE
  cart (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product (id),
    selected BOOLEAN NOT NULL DEFAULT TRUE,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    UNIQUE (user_id, product_id)
  );

CREATE TABLE
  checkout_info (
    id BIGSERIAL PRIMARY KEY,
    data JSONB NOT NULL,
    user_id BIGINT NOT NULL REFERENCES "user" (id)
  );