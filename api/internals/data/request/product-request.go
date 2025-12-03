package request

type CreateProductRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"required,min=1"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageUrl    string  `json:"image_url" binding:"required,url"`
	CategoryID  int64   `json:"category_id" binding:"required,min=1,max=100"`
	Discounted  *int32  `json:"discounted" binding:"omitempty,gte=0,lte=100"`
}

type UpdateProductRequest struct {
	Title       string  `json:"title" binding:"omitempty,min=1,max=255"`
	Description string  `json:"description" binding:"omitempty,min=1"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	ImageUrl    string  `json:"image_url" binding:"omitempty,url"`
	CategoryID  int64   `json:"category_id" binding:"omitempty,min=1,max=100"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
	Discounted  *int32  `json:"discounted" binding:"omitempty,gte=0,lte=100"`
}
