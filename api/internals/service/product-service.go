package service

import (
	"api/errors"
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
