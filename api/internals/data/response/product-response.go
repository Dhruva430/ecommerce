package response

import (
	"time"
)

type ProductResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	SellerID    int64     `json:"seller_id"`
	CreatedAt   time.Time `json:"created_at"`
	CategoryID  int64     `json:"category_id"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
