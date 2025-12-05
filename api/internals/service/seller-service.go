package service

import (
	"api/models/db"
	"context"
	"database/sql"
)

type SellerService struct {
	Queries *db.Queries
	Conn    *sql.DB
}

func NewSellerService(queries *db.Queries, conn *sql.DB) SellerService {
	return SellerService{
		Queries: queries,
		Conn:    conn,
	}
}

func SellerKYC(ctx context.Context, userID int64) error {
	// Implement KYC logic here
	return nil
}
