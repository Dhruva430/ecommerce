-- =========================================
-- ENUM TYPES
-- =========================================
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

-- =========================================
-- CORE USER TABLE
-- =========================================
CREATE TABLE
  "user" (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    role role NOT NULL DEFAULT 'BUYER',
    address_id BIGINT,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE
  );

-- =========================================
-- AUTH & PROFILE TABLES
-- =========================================
CREATE TABLE
  account (
    id BIGSERIAL PRIMARY KEY,
    account_id TEXT NOT NULL,
    provider provider NOT NULL DEFAULT 'CREDENTIALS',
    id_token TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    password TEXT,
    phone_number TEXT,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    UNIQUE (provider, user_id),
    UNIQUE (provider, account_id)
  );

CREATE TABLE
  refresh_token (
    id TEXT PRIMARY KEY,
    token TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    ip_address TEXT,
    last_used TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  buyer (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    loyalty_points INT NOT NULL DEFAULT 0,
    total_orders INT NOT NULL DEFAULT 0,
    wishlist JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  seller (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    status seller_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

CREATE TABLE
  seller_credentials (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT NOT NULL UNIQUE REFERENCES seller (id) ON DELETE CASCADE,
    business_name TEXT NOT NULL,
    gst_number TEXT,
    pan_number TEXT,
    bank_account_number TEXT,
    ifsc_code TEXT,
    business_address TEXT,
    website TEXT,
    contact_number BIGINT,
    contact_person TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
  );

CREATE TABLE
  admin (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

-- =========================================
-- ADDRESS
-- =========================================
CREATE TABLE
  address (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    pincode INT NOT NULL,
    area TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    country TEXT NOT NULL,
    phone_number BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    last_used TIMESTAMP NOT NULL DEFAULT NOW (),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
  );

ALTER TABLE "user" ADD CONSTRAINT fk_user_address FOREIGN KEY (address_id) REFERENCES address (id);

-- =========================================
-- PRODUCT CATEGORY
-- =========================================
CREATE TABLE
  product_category (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_by BIGINT,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    FOREIGN KEY (created_by) REFERENCES admin (id) ON DELETE SET NULL
  );

-- =========================================
-- PRODUCT
-- =========================================
CREATE TABLE
  product (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    seller_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    category_id BIGINT NOT NULL REFERENCES product_category (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id) ON DELETE CASCADE
  );

-- =========================================
-- PRODUCT VARIANTS + ATTRIBUTES + IMAGES
-- =========================================
CREATE TABLE
  product_variant (
    id BIGSERIAL PRIMARY KEY,
    size TEXT NOT NULL,
    description TEXT NOT NULL,
    discounted INT DEFAULT 0,
    title TEXT NOT NULL,
    price FLOAT NOT NULL,
    stock INT NOT NULL,
    product_id BIGINT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
  );

CREATE TABLE
  variant_images (
    id BIGSERIAL PRIMARY KEY,
    image_key TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    variant_id BIGINT NOT NULL,
    FOREIGN KEY (variant_id) REFERENCES product_variant (id) ON DELETE CASCADE
  );

CREATE TABLE
  variant_attribute (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    variant_id BIGINT NOT NULL,
    FOREIGN KEY (variant_id) REFERENCES product_variant (id) ON DELETE CASCADE
  );

-- =========================================
-- CART & REVIEWS & CHECKOUT
-- =========================================
CREATE TABLE
  cart (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    selected BOOLEAN NOT NULL DEFAULT TRUE,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product (id),
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
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE,
    UNIQUE (user_id, product_id)
  );

CREATE TABLE
  checkout_info (
    id BIGSERIAL PRIMARY KEY,
    data JSONB NOT NULL,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
  );

-- =========================================
-- ORDERS + ORDER PRODUCTS + PAYMENT
-- =========================================
CREATE TABLE
  orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    address JSONB NOT NULL,
    seller_id BIGINT,
    address_id BIGINT,
    total_amount FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    status order_status NOT NULL DEFAULT 'PENDING',
    payment_status payment_status NOT NULL DEFAULT 'PENDING',
    FOREIGN KEY (seller_id) REFERENCES seller (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (address_id) REFERENCES address (id)
  );

CREATE TABLE
  order_product (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,
    amount INT NOT NULL,
    variant_id BIGINT NOT NULL,
    order_id BIGINT,
    FOREIGN KEY (seller_id) REFERENCES seller (id),
    FOREIGN KEY (variant_id) REFERENCES product_variant (id),
    FOREIGN KEY (order_id) REFERENCES orders (id)
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
    FOREIGN KEY (order_id) REFERENCES orders (id)
  );

-- =========================================
-- SELLER DOCS & UPLOADS
-- =========================================
CREATE TABLE
  seller_documents (
    id BIGSERIAL PRIMARY KEY,
    document document_type NOT NULL DEFAULT 'IDENTITY_PROOF',
    document_url TEXT NOT NULL,
    seller_id BIGINT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW (),
    FOREIGN KEY (seller_id) REFERENCES seller (id) ON DELETE CASCADE
  );

CREATE TYPE upload_status AS ENUM ('PENDING', 'COMPLETED', 'Expired');

CREATE TABLE
  uploads (
    id BIGSERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    key TEXT NOT NULL,
    content_type TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    upload_type TEXT NOT NULL,
    status upload_status NOT NULL DEFAULT 'PENDING',
    user_id BIGINT NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW (),
    expires_at TIMESTAMP NOT NULL
  );

-- =========================================
-- VIEW
-- =========================================
CREATE
OR REPLACE VIEW user_view AS
SELECT
  *
FROM
  "user"
WHERE
  is_deleted = FALSE;

-- =========================================
-- INDEXES
-- =========================================
CREATE INDEX idx_cart_user ON cart (user_id);

CREATE INDEX idx_cart_product ON cart (product_id);

CREATE INDEX idx_order_user ON orders (user_id);

CREATE INDEX idx_order_seller ON orders (seller_id);

CREATE INDEX idx_product_seller ON product (seller_id);

CREATE INDEX idx_refresh_token_user ON refresh_token (user_id);