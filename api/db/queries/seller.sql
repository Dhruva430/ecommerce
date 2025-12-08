-- name: GetSellerByUserID :one
SELECT * FROM seller WHERE user_id = $1;

-- name: UpsertSellerCredentials :one
INSERT INTO seller_credentials (
  seller_id,
  business_name,
  gst_number,
  pan_number,
  bank_account_number,
  ifsc_code,
  business_address,
  website,
  contact_number,
  contact_person
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (seller_id) DO UPDATE SET
  business_name = EXCLUDED.business_name,
  gst_number = EXCLUDED.gst_number,
  pan_number = EXCLUDED.pan_number,
  bank_account_number = EXCLUDED.bank_account_number,
  ifsc_code = EXCLUDED.ifsc_code,
  business_address = EXCLUDED.business_address,
  website = EXCLUDED.website,
  contact_number = EXCLUDED.contact_number,
  contact_person = EXCLUDED.contact_person,
  updated_at = NOW()
RETURNING *;

-- name: CreateSellerDocument :one
INSERT INTO seller_documents (document, document_url, seller_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProductBySeller :one
SELECT * FROM product
WHERE id = $1 AND seller_id = $2;