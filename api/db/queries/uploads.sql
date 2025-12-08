-- name: CreateRequestFileUpload :one
INSERT INTO uploads (key, filename, content_type, file_size, upload_type, expires_at, user_id)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING *;

-- name: GetUploadRequestByKey :one
SELECT * FROM uploads WHERE key = $1;
