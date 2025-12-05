package request

type CreateProductRequest struct {
	Title       string                  `json:"title" binding:"required,min=1,max=255"`
	Description string                  `json:"description" binding:"required,min=1"`
	CategoryID  int64                   `json:"category_id" binding:"required,min=1,max=100"`
	Variant     []ProductVariantRequest `json:"variant" binding:"dive,required"`
}

type UpdateProductRequest struct {
	Title       string                  `json:"title" binding:"omitempty,min=1,max=255"`
	Description string                  `json:"description" binding:"omitempty,min=1"`
	CategoryID  int64                   `json:"category_id" binding:"omitempty,min=1,max=100"`
	IsActive    bool                    `json:"is_active" binding:"omitempty"`
	Variant     []ProductVariantRequest `json:"variant" binding:"dive,required"`
}

type ProductVariantRequest struct {
	Title       string                    `json:"title" binding:"required,min=1,max=255"`
	Description string                    `json:"description" binding:"required,min=1"`
	Size        string                    `json:"size" binding:"required,oneof=S M L XL XXL"`
	Discounted  int32                     `json:"discounted" binding:"omitempty,gte=0,lte=100"`
	Price       float64                   `json:"price" binding:"required,gt=0"`
	Stock       int32                     `json:"stock" binding:"required,gte=0"`
	Attributes  []VariantAttributeRequest `json:"attributes" binding:"dive,required"`
	Images      []VariantImageRequest     `json:"images" binding:"dive,required"`
}

type VariantAttributeRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Value string `json:"value" binding:"required,min=1,max=255"`
}

type VariantImageRequest struct {
	ImageKey string `json:"image_key" binding:"required"`
	Position int32  `json:"position" binding:"required,gte=0"`
}

// Example JSON
// {
//   "title": "Stylish Hoodie",
//   "description": "High-quality hoodie",
//   "category_id": 5,
//   "discounted": 10,
//   "variants": [
//     {
//       "title": "Red Hoodie",
//       "description": "Pure cotton red hoodie",
//       "size": "L",
//       "price": 1299,
//       "stock": 40,
//       "attributes": [
//         { "name": "color", "value": "red" },
//         { "name": "material", "value": "cotton" }
//       ],
//       "images": [
//         { "image_key": "product-images/red1.jpg", "position": 1 },
//         { "image_key": "product-images/red2.jpg", "position": 2 }
//       ]
//     }
//   ]
// }
