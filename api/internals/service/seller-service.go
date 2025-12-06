package service

import (
	"api/errors"
	"api/internals/data/request"
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

func (s *SellerService) SubmitKYC(ctx context.Context, userID int64, req request.SellerKYC) error {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		return errors.AppError{Message: "failed to get seller info", Code: 500}
	}
	if !seller.Verified {
		return errors.AppError{Message: "seller KYC not verified", Code: 403}
	}
	data := db.UpsertSellerCredentialsParams{
		SellerID:          seller.ID,
		BusinessName:      req.BusinessName,
		BusinessAddress:   sql.NullString{String: req.BusinessAddress, Valid: true},
		Website:           sql.NullString{String: req.Website, Valid: true},
		ContactNumber:     sql.NullInt64{Int64: req.ContactNumber, Valid: true},
		ContactPerson:     sql.NullString{String: req.ContactPerson, Valid: true},
		GstNumber:         sql.NullString{String: req.GstNumber, Valid: true},
		BankAccountNumber: sql.NullString{String: req.BankAccountNumber, Valid: true},
		IfscCode:          sql.NullString{String: req.IfscCode, Valid: true},
	}
	tx, err := s.Conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.AppError{Message: "failed to begin transaction", Code: 500}
	}
	qtx := s.Queries.WithTx(tx)
	_, err = qtx.UpsertSellerCredentials(ctx, data)
	if err != nil {
		tx.Rollback()
		return errors.AppError{Message: "failed to upsert seller credentials", Code: 500}
	}
	for _, doc := range req.Documents {
		docData := db.CreateSellerDocumentParams{
			Document:    db.DocumentType(doc.Document),
			DocumentUrl: doc.DocumentURL,
			SellerID:    seller.ID,
		}
		_, err = qtx.CreateSellerDocument(ctx, docData)
		if err != nil {
			tx.Rollback()
			return errors.AppError{Message: "failed to create seller document", Code: 500}
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.AppError{Message: "failed to commit transaction", Code: 500}
	}
	return nil
}
