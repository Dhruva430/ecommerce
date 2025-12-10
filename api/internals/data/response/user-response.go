package response

import (
	"database/sql"
	"time"
)

type UpdateAddressResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Pincode     int32     `json:"pin_code"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PhoneNumber int64     `json:"phone_number"`
	LastUsed    time.Time `json:"last_used"`
}
type GetUserAddressResponse struct {
	Addresses []UpdateAddressResponse `json:"addresses"`
}

type OrderResponse struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	AddressID     sql.NullInt64 `json:"address"`
	SellerID      int64         `json:"seller_id"`
	Total         float64       `json:"total_amount"`
	Status        string        `json:"status"`
	PaymentStatus string        `json:"payment_status"`
	CreatedAt     time.Time     `json:"created_at"`
}
type OrderHistoryResponse struct {
	Orders []OrderResponse `json:"orders"`
}
