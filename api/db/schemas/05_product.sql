CREATE TABLE
  product_category (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_by BIGINT REFERENCES admin (id) ON DELETE SET NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    created_at TIMESTAMP NOT NULL DEFAULT NOW ()
  );

CREATE TABLE
  product (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    seller_id BIGINT NOT NULL REFERENCES seller (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    category_id BIGINT NOT NULL REFERENCES product_category (id)
  );

CREATE TABLE
  product_variant (
    id BIGSERIAL PRIMARY KEY,
    size TEXT NOT NULL,
    description TEXT NOT NULL,
    discounted INT DEFAULT 0,
    title TEXT NOT NULL,
    price FLOAT NOT NULL,
    stock INT NOT NULL,
    product_id BIGINT NOT NULL REFERENCES product (id) ON DELETE CASCADE
  );

CREATE TABLE
  variant_images (
    id BIGSERIAL PRIMARY KEY,
    image_key TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    variant_id BIGINT NOT NULL REFERENCES product_variant (id) ON DELETE CASCADE
  );

CREATE TABLE
  variant_attribute (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    variant_id BIGINT NOT NULL REFERENCES product_variant (id) ON DELETE CASCADE
  );