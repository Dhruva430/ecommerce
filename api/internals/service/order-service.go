package service

import (
	"api/errors"
	"api/internals/data/request"
	"api/models/db"
	"context"
	"database/sql"
	"net/http"
)

type OrderService struct {
	Queries *db.Queries
	conn    *sql.DB
}

func NewOrderService(queries *db.Queries, conn *sql.DB) OrderService {
	return OrderService{
		Queries: queries,
		conn:    conn,
	}
}

func (s *OrderService) PlaceOrder(c context.Context, userID int64, req request.PlaceOrderRequest) error {
	tx, err := s.conn.BeginTx(c, nil)
	if err != nil {
		return &errors.AppError{Message: "failed to begin transaction", Code: http.StatusInternalServerError}
	}
	qtx := s.Queries.WithTx(tx)

	totalAmount := 0.0
	for _, item := range req.Items {
		variant, err := qtx.GetProductVariant(c, db.GetProductVariantParams{
			ProductID: item.ProductID,
			ID:        item.VariantID,
		})
		if err != nil {
			tx.Rollback()
			return &errors.AppError{Message: "failed to get product variant", Code: http.StatusInternalServerError}
		}
		if variant.Stock <= int32(item.Quantity) {
			tx.Rollback()
			return &errors.AppError{Message: "insufficient stock for product variant", Code: http.StatusBadRequest}
		}

		totalAmount += float64(item.Quantity) * variant.Price

	}

	order, err := qtx.CreateOrder(c, db.CreateOrderParams{
		UserID:      userID,
		AddressID:   sql.NullInt64{Int64: req.AddressID, Valid: req.AddressID > 0},
		TotalAmount: totalAmount,
	})
	if err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to create order", Code: http.StatusInternalServerError}
	}
	for _, item := range req.Items {

		_, err := qtx.CreateOrderProduct(c, db.CreateOrderProductParams{
			OrderID:   sql.NullInt64{Int64: order.ID, Valid: order.ID > 0},
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Amount:    int32(item.Quantity),
		})
		if err != nil {
			tx.Rollback()
			return &errors.AppError{Message: "failed to create order items", Code: http.StatusInternalServerError}
		}

		err = qtx.DecrementProductVariantStock(c, db.DecrementProductVariantStockParams{
			ID:    item.VariantID,
			Stock: int32(item.Quantity),
		})
		if err != nil {
			tx.Rollback()
			return &errors.AppError{Message: "failed to reduce stock", Code: http.StatusInternalServerError}
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}

	return nil
}

func (s *OrderService) GetOrderDetails(c context.Context, userID int64, orderID int64) (db.Order, error) {
	order, err := s.Queries.GetOrderByID(c, db.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		return db.Order{}, &errors.AppError{Message: "failed to get order details", Code: http.StatusInternalServerError}
	}
	return order, nil
}
