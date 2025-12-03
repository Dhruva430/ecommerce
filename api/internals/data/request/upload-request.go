package request

type RequestFileUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required"`
	UploadType  string `json:"upload_type" binding:"required,oneof=product-image"`
}
