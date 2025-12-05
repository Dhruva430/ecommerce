package request

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof= SELLER BUYER"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
