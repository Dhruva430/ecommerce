package middleware

import (
	"api/errors"
	"api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := util.ExtractTokenFromHeader(c)
		if err != nil {
			c.Error(&errors.AppError{Message: "authorization token required", Code: http.StatusUnauthorized})
			return
		}

		claims, err := util.ParseJWT(tokenString)
		if err != nil {
			c.Error(&errors.AppError{Message: "invalid token", Code: http.StatusUnauthorized})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
