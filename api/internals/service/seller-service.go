package service

import (
	"api/errors"
	"api/internals/data/request"
	"api/internals/data/response"
	"api/models/db"
	"api/util"
	"context"
	"database/sql"
	"net/http"
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

func (s *SellerService) ApplyForKYC(ctx context.Context, userID int64, req request.SellerKYC) error {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		return errors.AppError{Message: "failed to get seller info", Code: 500}
	}
	if seller.ID == 0 {
		return errors.AppError{Message: "seller not found", Code: 404}
	}
	if seller.Verified {
		return errors.AppError{Message: "seller already verified", Code: 400}
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
func (s *SellerService) CreateProduct(ctx context.Context, req request.CreateProductRequest, userID int64) (response.ProductResponse, error) {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ProductResponse{}, &errors.AppError{Message: "user is not a seller", Code: http.StatusForbidden}
		}
		return response.ProductResponse{}, &errors.AppError{Message: "failed to verify seller", Code: http.StatusInternalServerError}
	}

	if seller.Status != db.SellerStatusAPPROVED {
		return response.ProductResponse{}, &errors.AppError{Message: "seller is not approved", Code: http.StatusForbidden}
	}

	tx, err := s.Conn.BeginTx(ctx, nil)
	if err != nil {
		return response.ProductResponse{}, &errors.AppError{Message: "failed to start transaction", Code: http.StatusInternalServerError}
	}
	qtx := s.Queries.WithTx(tx)

	productData := db.CreateProductParams{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		SellerID:    seller.ID,
	}

	product, err := qtx.CreateProduct(ctx, productData)
	if err != nil {
		tx.Rollback()
		return response.ProductResponse{}, &errors.AppError{Message: "failed to create product", Code: http.StatusInternalServerError}
	}

	for _, variant := range req.Variant {
		variantData := db.CreateProductVariantParams{
			ProductID:   product.ID,
			Title:       variant.Title,
			Description: variant.Description,
			Size:        variant.Size,
			Price:       variant.Price,
			Stock:       variant.Stock,
		}

		v, err := qtx.CreateProductVariant(ctx, variantData)
		if err != nil {
			tx.Rollback()
			return response.ProductResponse{}, &errors.AppError{Message: "failed to create product variant", Code: http.StatusInternalServerError}
		}
		for _, attr := range variant.Attributes {
			attrData := db.CreateVariantAttributeParams{
				VariantID: v.ID,
				Name:      attr.Name,
				Value:     attr.Value,
			}
			_, err := qtx.CreateVariantAttribute(ctx, attrData)
			if err != nil {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "failed to create variant attribute", Code: http.StatusInternalServerError}
			}
		}
		for _, img := range variant.Images {
			uploadReq, err := qtx.GetUploadRequestByKey(ctx, img.ImageKey)
			if err != nil {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "failed to get upload request", Code: http.StatusInternalServerError}
			}
			if uploadReq.UserID != seller.ID {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "upload request does not belong to seller", Code: http.StatusForbidden}
			}
			if uploadReq.Status != "COMPLETED" {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "image upload not completed", Code: http.StatusBadRequest}
			}
			_, err = qtx.CreateVariantImage(ctx, db.CreateVariantImageParams{
				VariantID: v.ID,
				ImageKey:  img.ImageKey,
				Position:  img.Position,
			})
			if err != nil {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "failed to create variant image", Code: http.StatusInternalServerError}
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return response.ProductResponse{}, &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}

	return response.ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		IsActive:    product.IsActive,
		SellerID:    product.SellerID,
		CreatedAt:   product.CreatedAt,
		CategoryID:  product.CategoryID,
	}, nil
}
func (s *SellerService) UpdateProduct(ctx context.Context, productID int64, req request.UpdateProductRequest, userID int64) (response.ProductResponse, error) {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ProductResponse{}, &errors.AppError{Message: "user is not a seller", Code: http.StatusForbidden}
		}
		return response.ProductResponse{}, &errors.AppError{Message: "failed to verify seller", Code: http.StatusInternalServerError}
	}

	product, err := s.Queries.GetProductBySeller(ctx, db.GetProductBySellerParams{
		ID:       productID,
		SellerID: seller.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ProductResponse{}, &errors.AppError{Message: "product not found or you don't have permission", Code: http.StatusNotFound}
		}
		return response.ProductResponse{}, &errors.AppError{Message: "failed to get product", Code: http.StatusInternalServerError}
	}

	tx, err := s.Conn.BeginTx(ctx, nil)
	if err != nil {
		return response.ProductResponse{}, &errors.AppError{Message: "failed to start transaction", Code: http.StatusInternalServerError}
	}
	qtx := s.Queries.WithTx(tx)

	productData := db.UpdateProductParams{
		ID:          product.ID,
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		IsActive:    req.IsActive,
	}
	product, err = qtx.UpdateProduct(ctx, productData)
	if err != nil {
		tx.Rollback()
		return response.ProductResponse{}, &errors.AppError{Message: "failed to update product", Code: http.StatusInternalServerError}
	}
	for _, variant := range req.Variant {

		variantData := db.UpdateProductVariantParams{
			ID:          product.ID,
			Title:       variant.Title,
			Description: variant.Description,
			Size:        variant.Size,
			Price:       variant.Price,
			Discounted:  sql.NullInt32{Int32: variant.Discounted, Valid: variant.Discounted >= 0},
			Stock:       variant.Stock,
		}

		v, err := qtx.UpdateProductVariant(ctx, variantData)
		if err != nil {
			tx.Rollback()
			return response.ProductResponse{}, &errors.AppError{Message: "failed to update product variant", Code: http.StatusInternalServerError}
		}
		for _, attr := range variant.Attributes {
			attrData := db.UpdateVariantAttributeParams{
				ID:    v.ID,
				Name:  attr.Name,
				Value: attr.Value,
			}
			_, err := qtx.UpdateVariantAttribute(ctx, attrData)
			if err != nil {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "failed to update variant attribute", Code: http.StatusInternalServerError}
			}
		}
		for _, img := range variant.Images {
			uploadReq, err := qtx.GetUploadRequestByKey(ctx, img.ImageKey)
			if err != nil {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "failed to get upload request", Code: http.StatusInternalServerError}
			}
			if uploadReq.UserID != seller.ID {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "upload request does not belong to seller", Code: http.StatusForbidden}
			}
			if uploadReq.Status != "COMPLETED" {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "image upload not completed", Code: http.StatusBadRequest}
			}
			_, err = qtx.UpdateVariantImage(ctx, db.UpdateVariantImageParams{
				ID:       v.ID,
				ImageKey: img.ImageKey,
				Position: img.Position,
			})
			if err != nil {
				tx.Rollback()
				return response.ProductResponse{}, &errors.AppError{Message: "failed to update variant image", Code: http.StatusInternalServerError}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return response.ProductResponse{}, &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}

	return response.ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		IsActive:    product.IsActive,
		SellerID:    product.SellerID,
		CreatedAt:   product.CreatedAt,
		CategoryID:  product.CategoryID,
	}, nil

}

func (s *SellerService) DeleteProduct(ctx context.Context, productID int64, userID int64) error {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &errors.AppError{Message: "user is not a seller", Code: http.StatusForbidden}
		}
		return &errors.AppError{Message: "failed to verify seller", Code: http.StatusInternalServerError}
	}

	_, err = s.Queries.GetProductBySeller(ctx, db.GetProductBySellerParams{
		ID:       productID,
		SellerID: seller.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &errors.AppError{Message: "product not found or you don't have permission", Code: http.StatusNotFound}
		}
		return &errors.AppError{Message: "failed to get product", Code: http.StatusInternalServerError}
	}

	err = s.Queries.DeleteProduct(ctx, productID)
	if err != nil {
		return &errors.AppError{Message: "failed to delete product", Code: http.StatusInternalServerError}
	}

	return nil
}

func (s *SellerService) RegisterSeller(ctx context.Context, req request.RegisterSellerRequest) error {
	existingUser, err := s.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return &errors.AppError{Message: "failed to check existing user", Code: http.StatusInternalServerError}
	}
	if existingUser.ID != 0 {
		return &errors.AppError{Message: "user with this email already exists", Code: http.StatusBadRequest}
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &errors.AppError{Message: "failed to hash password", Code: http.StatusInternalServerError}
	}
	userData := db.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
	}
	tx, err := s.Conn.BeginTx(ctx, nil)
	if err != nil {
		return &errors.AppError{Message: "failed to begin transaction", Code: http.StatusInternalServerError}
	}
	qtx := s.Queries.WithTx(tx)
	user, err := qtx.CreateUser(ctx, userData)
	if err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to create user", Code: http.StatusInternalServerError}
	}
	accountData := db.CreateAccountParams{
		UserID:    user.ID,
		Provider:  db.ProviderCREDENTIALS,
		AccountID: req.Email,
		Password:  sql.NullString{String: hashedPassword, Valid: true},
	}
	if err = qtx.CreateAccount(ctx, accountData); err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to create account", Code: http.StatusInternalServerError}
	}
	_, err = qtx.CreateSeller(ctx, user.ID)
	if err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to create seller", Code: http.StatusInternalServerError}
	}

	err = tx.Commit()
	if err != nil {
		return &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}
	return nil

}
func (s *SellerService) LoginSeller(ctx context.Context, req request.LoginRequest, ip string) error {
	return nil
}

// TODO: Add analytics methods for seller service
