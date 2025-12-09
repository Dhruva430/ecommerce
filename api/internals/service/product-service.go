package service

import (
	"api/errors"
	"api/internals/data/request"
	"api/internals/data/response"
	"api/models/db"
	"context"
	"database/sql"
	"net/http"
	"strconv"
)

type ProductService struct {
	Queries *db.Queries
	Conn    *sql.DB
}

func NewProductService(queries *db.Queries, conn *sql.DB) ProductService {
	return ProductService{
		Queries: queries,
		Conn:    conn,
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context, page, pageSize int, category, sellerID, isActive string) ([]response.ProductResponse, int64, error) {
	offset := (page - 1) * pageSize

	var categoryFilter sql.NullString
	if category != "" {
		categoryFilter = sql.NullString{String: category, Valid: true}
	}

	var sellerIDFilter sql.NullInt64
	if sellerID != "" {
		id, err := strconv.ParseInt(sellerID, 10, 64)
		if err == nil {
			sellerIDFilter = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	var isActiveFilter sql.NullBool
	if isActive != "" {
		active, err := strconv.ParseBool(isActive)
		if err == nil {
			isActiveFilter = sql.NullBool{Bool: active, Valid: true}
		}
	}

	products, err := s.Queries.GetAllProducts(ctx, db.GetAllProductsParams{
		CategoryID:  categoryFilter,
		SellerID:    sellerIDFilter,
		IsActive:    isActiveFilter,
		LimitCount:  int32(pageSize),
		OffsetCount: int32(offset),
	})
	if err != nil {
		return nil, 0, &errors.AppError{Message: "failed to get products", Code: http.StatusInternalServerError}
	}

	count, err := s.Queries.CountProducts(ctx, db.CountProductsParams{
		CategoryID: categoryFilter,
		SellerID:   sellerIDFilter,
		IsActive:   isActiveFilter,
	})
	if err != nil {
		return nil, 0, &errors.AppError{Message: "failed to count products", Code: http.StatusInternalServerError}
	}

	var respProducts []response.ProductResponse
	for _, product := range products {

		respProducts = append(respProducts, response.ProductResponse{
			ID:          product.ID,
			Title:       product.Title,
			Description: product.Description,
			IsActive:    product.IsActive,
			SellerID:    product.SellerID,
			CreatedAt:   product.CreatedAt,
			CategoryID:  product.CategoryID,
		})
	}

	return respProducts, count, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (response.ProductResponse, error) {
	product, err := s.Queries.GetProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ProductResponse{}, &errors.AppError{Message: "product not found", Code: http.StatusNotFound}
		}
		return response.ProductResponse{}, &errors.AppError{Message: "failed to get product", Code: http.StatusInternalServerError}
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

func (s *ProductService) AddProductVariant(ctx context.Context, userID int64, productID int64, req []request.ProductVariantRequest) error {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		return &errors.AppError{Message: "failed to get seller", Code: http.StatusInternalServerError}
	}

	if seller.ID == 0 {
		return &errors.AppError{Message: "seller not found", Code: http.StatusNotFound}
	}
	if seller.Verified {
		return &errors.AppError{Message: "seller KYC not verified", Code: http.StatusForbidden}
	}
	product, err := s.Queries.GetProductByID(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &errors.AppError{Message: "product not found", Code: http.StatusNotFound}
		}
		return &errors.AppError{Message: "failed to get product", Code: http.StatusInternalServerError}
	}

	if product.SellerID != seller.ID {
		return &errors.AppError{Message: "unauthorized to add variant to this product", Code: http.StatusUnauthorized}
	}
	tx, err := s.Conn.BeginTx(ctx, nil)
	if err != nil {
		return &errors.AppError{Message: "failed to start transaction", Code: http.StatusInternalServerError}
	}
	qtx := s.Queries.WithTx(tx)

	for _, variant := range req {
		v, err := qtx.CreateProductVariant(ctx, db.CreateProductVariantParams{
			ProductID:   productID,
			Title:       variant.Title,
			Description: variant.Description,
			Size:        variant.Size,
			Price:       variant.Price,
			Discounted:  sql.NullInt32{Int32: int32(variant.Discounted), Valid: variant.Discounted > 0},
			Stock:       int32(variant.Stock),
		})
		if err != nil {
			tx.Rollback()
			return &errors.AppError{Message: "failed to add product variant", Code: http.StatusInternalServerError}
		}
		for _, attr := range variant.Attributes {
			_, err := qtx.CreateVariantAttribute(ctx, db.CreateVariantAttributeParams{
				VariantID: v.ID,
				Name:      attr.Name,
				Value:     attr.Value,
			})
			if err != nil {
				tx.Rollback()
				return &errors.AppError{Message: "failed to add variant attribute", Code: http.StatusInternalServerError}
			}
		}
		for _, img := range variant.Images {
			uploadReq, err := qtx.GetUploadRequestByKey(ctx, img.ImageKey)
			if err != nil {
				tx.Rollback()
				return &errors.AppError{Message: "failed to get upload request", Code: http.StatusInternalServerError}
			}
			if uploadReq.ID != seller.ID {
				tx.Rollback()
				return &errors.AppError{Message: "image does not belong to seller", Code: http.StatusBadRequest}
			}

			if uploadReq.Status != "COMPLETED" {
				tx.Rollback()
				return &errors.AppError{Message: "image upload not completed", Code: http.StatusBadRequest}
			}
			_, err = qtx.CreateVariantImage(ctx, db.CreateVariantImageParams{
				VariantID: v.ID,
				ImageKey:  img.ImageKey,
				Position:  int32(img.Position),
			})
			if err != nil {
				tx.Rollback()
				return &errors.AppError{Message: "failed to add variant image", Code: http.StatusInternalServerError}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}
	return nil

}
