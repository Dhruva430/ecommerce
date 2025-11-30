package request

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof= SELLER BUYER"`
	Password string `json:"password" binding:"required,min=8"`
}
