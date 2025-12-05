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
		ID: product.ID,
	}, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, productID int64, req request.UpdateProductRequest, userID int64) (response.ProductResponse, error) {
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

func (s *ProductService) DeleteProduct(ctx context.Context, productID int64, userID int64) error {
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
