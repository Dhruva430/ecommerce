package request

type PlaceOrderRequest struct {
	Items         []OrderItemRequest `json:"items" binding:"required,dive,required"`
	AddressID     int64              `json:"address_id" binding:"required"`
	UserID        int64              `json:"user_id" binding:"-"`
	PaymentMethod string             `json:"payment_method" binding:"required,oneof=credit_card debit_card paypal cod"`
}

type OrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	VariantID int64 `json:"variant_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending shipped delivered canceled"`
}
