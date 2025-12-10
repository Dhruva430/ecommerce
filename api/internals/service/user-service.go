package service

import (
	"api/internals/data/request"
	"api/internals/data/response"
	"api/models/db"
	"context"
	"database/sql"
)

type UserService struct {
	Queries *db.Queries
	Conn    *sql.DB
}

func NewUserService(queries *db.Queries, conn *sql.DB) UserService {
	return UserService{
		Queries: queries,
		Conn:    conn,
	}
}

func (s *UserService) GetUserAddress(userID int64) (response.GetUserAddressResponse, error) {
	addresses, err := s.Queries.GetUserAddresses(nil, userID)
	if err != nil {
		return response.GetUserAddressResponse{}, err
	}
	if len(addresses) == 0 {
		return response.GetUserAddressResponse{}, sql.ErrNoRows
	}
	var respAddresses []response.UpdateAddressResponse
	for _, addr := range addresses {
		respAddr := response.UpdateAddressResponse{
			ID:          addr.ID,
			Name:        addr.Name,
			Pincode:     addr.Pincode,
			City:        addr.City,
			State:       addr.State,
			Country:     addr.Country,
			PhoneNumber: addr.PhoneNumber,
			LastUsed:    addr.LastUsed,
		}
		respAddresses = append(respAddresses, respAddr)
	}
	return response.GetUserAddressResponse{Addresses: respAddresses}, nil
}
func (s *UserService) UpdateUserAddress(c context.Context, req request.UpdateAddressRequest) error {

	addressParams := db.UpdateUserAddressParams{
		ID:          req.AddressID,
		Name:        req.Name,
		Pincode:     req.PinCode,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		PhoneNumber: req.PhoneNumber,
	}
	if err := s.Queries.UpdateUserAddress(c, addressParams); err != nil {
		return err
	}
	return nil

}
func (s *UserService) GetOrderHistory(ctx context.Context, userID int64, filter string) (response.OrderHistoryResponse, error) {
	orders, err := s.Queries.GetOrderHistory(ctx, userID)
	if err != nil {
		return response.OrderHistoryResponse{}, err
	}

	var respOrders []response.OrderResponse
	for _, order := range orders {
		respOrder := response.OrderResponse{
			ID:            order.ID,
			UserID:        order.UserID,
			AddressID:     order.AddressID,
			Total:         order.TotalAmount,
			Status:        string(order.Status),
			PaymentStatus: string(order.PaymentStatus),
			CreatedAt:     order.CreatedAt,
		}
		respOrders = append(respOrders, respOrder)
	}

	return response.OrderHistoryResponse{Orders: respOrders}, nil
}
func (s *UserService) DeleteUser(c context.Context, userID int64) error {
	err := s.Queries.DeleteUser(c, userID)
	if err != nil {
		return err
	}
	return nil
}
