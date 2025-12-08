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
    user_id BIGINT NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    UNIQUE (provider, user_id),
    UNIQUE (provider, account_id)
  );

CREATE TABLE
  refresh_token (
    id TEXT PRIMARY KEY,
    token TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    ip_address TEXT,
    last_used TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    expires_at TIMESTAMP NOT NULL
  );