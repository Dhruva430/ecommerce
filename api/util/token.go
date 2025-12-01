package util

import (
	"api/configs"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID      int32   `json:"user_id"`
	Type        string  `json:"type"`
	Permissions []int32 `json:"perms,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(subject int32, audience string, perms []int32) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := JWTClaims{
		UserID:      subject,
		Type:        "access",
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := configs.GetJWTSecret()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

func GenerateRefreshToken(subject int32, audience string) (string, *JWTClaims, error) {
	expiration := time.Now().Add(7 * 24 * time.Hour)
	claims := JWTClaims{
		UserID: subject,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := configs.GetJWTSecret()
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, &claims, err
}

func ParseJWT(tokenString string) (*JWTClaims, error) {
	secret := configs.GetJWTSecret()
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.Split(bearer, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
