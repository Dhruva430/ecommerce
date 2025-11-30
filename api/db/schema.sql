CREATE TYPE seller_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TYPE role AS ENUM ('BUYER', 'SELLER', 'ADMIN');

CREATE TYPE order_status AS ENUM (
  'PENDING',
  'PROCESSING',
  'SHIPPED',
  'DELIVERED',
  'CANCELLED'
);

CREATE TYPE payment_status AS ENUM ('PENDING', 'COMPLETED', 'FAILED', 'REFUNDED');

CREATE TYPE document_type AS ENUM (
  'GST_CERTIFICATE',
  'BUSINESS_LICENSE',
  'IDENTITY_PROOF'
);

CREATE TYPE provider AS ENUM ('GOOGLE', 'FACEBOOK', 'GITHUB', 'CREDENTIALS');

CREATE INDEX idx_cart_user ON cart (user_id);

CREATE INDEX idx_cart_product ON cart (product_id);

CREATE INDEX idx_order_user ON "order" (user_id);

CREATE INDEX idx_order_seller ON "order" (seller_id);

CREATE INDEX idx_product_seller ON product (seller_id);

CREATE INDEX idx_refresh_token_user ON refresh_token (user_id);

CREATE TABLE
  "user" (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    role role NOT NULL DEFAULT 'BUYER',
    address_id BIGINT,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user_address FOREIGN KEY (address_id) REFERENCES address (id)
  );

CREATE TABLE
  account (
    id BIGSERIAL PRIMARY KEY,
    account_id TEXT NOT NULL UNIQUE,
    provider provider NOT NULL DEFAULT 'CREDENTIALS',
    id_token TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    password TEXT,
    phone_number TEXT,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_creds_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    UNIQUE (provider_id, account_id)
  );

CREATE TABLE
  refresh_token (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    ip_address TEXT,
    device TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_refresh_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  buyer (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    loyalty_points INT NOT NULL DEFAULT 0,
    total_orders INT NOT NULL DEFAULT 0,
    wishlist JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_buyer_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  seller (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    shop_name TEXT,
    shop_logo TEXT,
    description TEXT,
    gst_number TEXT,
    status seller_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_seller_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  admin (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_admin_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  address (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    pincode INT NOT NULL,
    house_no TEXT NOT NULL,
    area TEXT NOT NULL,
    landmark TEXT,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    phone_number BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    last_used TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_address_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  product (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price FLOAT NOT NULL,
    image_url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    discounted INT DEFAULT 0,
    seller_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    category TEXT NOT NULL,
    CONSTRAINT fk_product_seller FOREIGN KEY (seller_id) REFERENCES seller (id) ON DELETE CASCADE
  );

CREATE TABLE
  cart (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    selected BOOLEAN NOT NULL DEFAULT TRUE,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_cart_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_product FOREIGN KEY (product_id) REFERENCES product (id),
    UNIQUE (user_id, product_id)
  );

CREATE TABLE
  reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    rating INT NOT NULL,
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_product FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE,
    UNIQUE (user_id, product_id)
  );

CREATE TABLE
  checkout_info (
    id BIGSERIAL PRIMARY KEY,
    data JSONB NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_checkout_user FOREIGN KEY (user_id) REFERENCES "user" (id)
  );

CREATE TABLE
  "order" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    address JSONB NOT NULL,
    seller_id BIGINT,
    address_id BIGINT,
    total_amount FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    status order_status NOT NULL DEFAULT 'PENDING',
    payment_status payment_status NOT NULL DEFAULT 'PENDING',
    CONSTRAINT fk_order_seller FOREIGN KEY (seller_id) REFERENCES seller (id),
    CONSTRAINT fk_order_user FOREIGN KEY (user_id) REFERENCES "user" (id),
    CONSTRAINT fk_order_address FOREIGN KEY (address_id) REFERENCES address (id)
  );

CREATE TABLE
  product_variant (
    id BIGSERIAL PRIMARY KEY,
    size TEXT NOT NULL,
    description TEXT NOT NULL,
    title TEXT NOT NULL,
    price FLOAT NOT NULL,
    stock INT NOT NULL,
    product_id BIGINT NOT NULL,
    CONSTRAINT fk_variant_product FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
  );

CREATE TABLE
  order_product (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,
    amount INT NOT NULL,
    variant_id BIGINT NOT NULL,
    order_id BIGINT,
    CONSTRAINT fk_op_seller FOREIGN KEY (seller_id) REFERENCES seller (id),
    CONSTRAINT fk_op_variant FOREIGN KEY (variant_id) REFERENCES product_variant (id),
    CONSTRAINT fk_op_order FOREIGN KEY (order_id) REFERENCES "order" (id)
  );

CREATE TABLE
  payment (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE,
    transaction_id TEXT UNIQUE,
    payment_gateway TEXT,
    amount FLOAT NOT NULL,
    status payment_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES "order" (id)
  );

CREATE TABLE
  variant_images (
    id BIGSERIAL PRIMARY KEY,
    image_url TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    variant_id BIGINT NOT NULL,
    CONSTRAINT fk_varimg_variant FOREIGN KEY (variant_id) REFERENCES product_variant (id) ON DELETE CASCADE
  );

CREATE TABLE
  variant_attribute (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    variant_id BIGINT NOT NULL,
    CONSTRAINT fk_varattr_variant FOREIGN KEY (variant_id) REFERENCES product_variant (id) ON DELETE CASCADE
  );

CREATE TABLE
  seller_documents (
    id BIGSERIAL PRIMARY KEY,
    document document_type NOT NULL DEFAULT 'IDENTITY_PROOF',
    document_url TEXT NOT NULL,
    seller_id BIGINT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW (),
    CONSTRAINT fk_sellerdoc_seller FOREIGN KEY (seller_id) REFERENCES seller (id) ON DELETE CASCADE
  );