package service

import (
	"api/errors"
	awsclient "api/internals/aws"
	"api/internals/data/request"
	"api/internals/data/response"
	"api/models/db"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const MAX_FILE_SIZE = 10 * 1024 * 1024

type UploadService struct {
	Queries *db.Queries
	Conn    *sql.DB
}

func NewUploadService(queries *db.Queries, conn *sql.DB) UploadService {
	return UploadService{
		Queries: queries,
		Conn:    conn,
	}
}

func generateFileKey() string {
	return uuid.New().String()
}
func (u *UploadService) RequestFileUpload(ctx context.Context, userID int64, req request.RequestFileUploadRequest) (response.RequestFileUploadResponse, error) {
	if req.FileSize > MAX_FILE_SIZE {
		return response.RequestFileUploadResponse{}, errors.AppError{Message: "file size exceeds limit", Code: 400}
	}
	fileKey := generateFileKey()
	uploadFolder := ""
	switch req.UploadType {
	case "product-image":
		uploadFolder = "product-images/"
	default:
		return response.RequestFileUploadResponse{}, errors.AppError{Message: "invalid upload type", Code: 400}

	}
	data := db.CreateRequestFileUploadParams{
		Key:         uploadFolder + fileKey,
		Filename:    req.Filename,
		ContentType: req.ContentType,
		FileSize:    req.FileSize,
		UploadType:  req.UploadType,
		ExpiresAt:   time.Now().Add(15 * time.Minute),
	}
	_, err := u.Queries.CreateRequestFileUpload(ctx, data)
	if err != nil {
		return response.RequestFileUploadResponse{}, errors.AppError{Message: "failed to create upload request", Code: 500}
	}
	presignedurl, err := awsclient.GeneratePresignedUploadURL(data.Key, data.ContentType, data.FileSize)
	if err != nil {
		return response.RequestFileUploadResponse{}, errors.AppError{Message: "failed to generate presigned URL", Code: 500}
	}
	return response.RequestFileUploadResponse{
		FileKey:     data.Key,
		UploadURL:   presignedurl,
		ExpiresIn:   int64(time.Until(data.ExpiresAt)),
		UploadType:  data.UploadType,
		ContentType: data.ContentType,
		FileSize:    data.FileSize,
		UserID:      userID,
		Bucket:      "ecommerce-dhruva",
	}, nil

}
