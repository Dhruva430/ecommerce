package response

type RequestFileUploadResponse struct {
	FileKey     string `json:"file_key"`
	UploadURL   string `json:"upload_url"`
	ExpiresIn   int64  `json:"expires_in"`
	UploadType  string `json:"upload_type"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
	UserID      int64  `json:"user_id"`
	Bucket      string `json:"bucket"`
}
