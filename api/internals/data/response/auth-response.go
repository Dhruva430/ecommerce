package response

import (
	"time"
)

type UserResponse struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username,omitempty"`
	AccessToken    string    `json:"access_token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
	TokenExpiresAt time.Time `json:"token_expires_at,omitempty"`
}

type TokenResponse struct {
	RefreshToken   string    `json:"refresh_token"`
	AccessToken    string    `json:"access_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
}

type UserDetails struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username,omitempty"`
}
