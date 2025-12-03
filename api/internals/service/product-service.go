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
		var discounted *int32
		if product.Discounted.Valid {
			discounted = &product.Discounted.Int32
		}

		respProducts = append(respProducts, response.ProductResponse{
			ID:          product.ID,
			Title:       product.Title,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			IsActive:    product.IsActive,
			Discounted:  discounted,
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

	var discounted *int32
	if product.Discounted.Valid {
		discounted = &product.Discounted.Int32
	}

	return response.ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		IsActive:    product.IsActive,
		Discounted:  discounted,
		SellerID:    product.SellerID,
		CreatedAt:   product.CreatedAt,
		CategoryID:  product.CategoryID,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req request.CreateProductRequest, userID int64) (response.ProductResponse, error) {
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

	var discounted sql.NullInt32
	if req.Discounted != nil {
		discounted = sql.NullInt32{Int32: *req.Discounted, Valid: true}
	}

	product, err := s.Queries.CreateProduct(ctx, db.CreateProductParams{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		ImageUrl:    req.ImageUrl,
		CategoryID:  req.CategoryID,
		Discounted:  discounted,
		SellerID:    seller.ID,
	})
	if err != nil {
		return response.ProductResponse{}, &errors.AppError{Message: "failed to create product", Code: http.StatusInternalServerError}
	}

	var respDiscounted *int32
	if product.Discounted.Valid {
		respDiscounted = &product.Discounted.Int32
	}

	return response.ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		IsActive:    product.IsActive,
		Discounted:  respDiscounted,
		SellerID:    product.SellerID,
		CreatedAt:   product.CreatedAt,
		CategoryID:  product.CategoryID,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req request.UpdateProductRequest, userID int64) (response.ProductResponse, error) {
	// TODO: Change to Put method for full update
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ProductResponse{}, &errors.AppError{Message: "user is not a seller", Code: http.StatusForbidden}
		}
		return response.ProductResponse{}, &errors.AppError{Message: "failed to verify seller", Code: http.StatusInternalServerError}
	}

	existingProduct, err := s.Queries.GetProductBySeller(ctx, db.GetProductBySellerParams{
		ID:       id,
		SellerID: seller.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ProductResponse{}, &errors.AppError{Message: "product not found or you don't have permission", Code: http.StatusNotFound}
		}
		return response.ProductResponse{}, &errors.AppError{Message: "failed to get product", Code: http.StatusInternalServerError}
	}

	var title, description, imageUrl sql.NullString
	var price sql.NullFloat64
	var isActive sql.NullBool
	var discounted sql.NullInt32
	var category sql.NullInt64
	if req.Title != "" {
		title = sql.NullString{String: req.Title, Valid: true}
	}
	if req.Description != "" {
		description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.Price > 0 {
		price = sql.NullFloat64{Float64: req.Price, Valid: true}
	}
	if req.ImageUrl != "" {
		imageUrl = sql.NullString{String: req.ImageUrl, Valid: true}
	}
	if req.CategoryID != 0 {
		category = sql.NullInt64{Int64: req.CategoryID, Valid: true}
	}
	if req.IsActive != nil {
		isActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	}
	if req.Discounted != nil {
		discounted = sql.NullInt32{Int32: *req.Discounted, Valid: true}
	}

	product, err := s.Queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          existingProduct.ID,
		Title:       title,
		Description: description,
		Price:       price,
		ImageUrl:    imageUrl,
		CategoryID:  category,
		IsActive:    isActive,
		Discounted:  discounted,
	})
	if err != nil {
		return response.ProductResponse{}, &errors.AppError{Message: "failed to update product", Code: http.StatusInternalServerError}
	}

	var respDiscounted *int32
	if product.Discounted.Valid {
		respDiscounted = &product.Discounted.Int32
	}

	return response.ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		IsActive:    product.IsActive,
		Discounted:  respDiscounted,
		SellerID:    product.SellerID,
		CreatedAt:   product.CreatedAt,
		CategoryID:  product.CategoryID,
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64, userID int64) error {
	seller, err := s.Queries.GetSellerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &errors.AppError{Message: "user is not a seller", Code: http.StatusForbidden}
		}
		return &errors.AppError{Message: "failed to verify seller", Code: http.StatusInternalServerError}
	}

	_, err = s.Queries.GetProductBySeller(ctx, db.GetProductBySellerParams{
		ID:       id,
		SellerID: seller.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &errors.AppError{Message: "product not found or you don't have permission", Code: http.StatusNotFound}
		}
		return &errors.AppError{Message: "failed to get product", Code: http.StatusInternalServerError}
	}

	err = s.Queries.DeleteProduct(ctx, id)
	if err != nil {
		return &errors.AppError{Message: "failed to delete product", Code: http.StatusInternalServerError}
	}

	return nil
}
