CREATE TABLE
  seller_documents (
    id BIGSERIAL PRIMARY KEY,
    document document_type NOT NULL DEFAULT 'IDENTITY_PROOF',
    document_url TEXT NOT NULL,
    seller_id BIGINT NOT NULL REFERENCES seller (id) ON DELETE CASCADE,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW ()
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