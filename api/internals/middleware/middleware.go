package middleware

import (
	"api/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	fmt.Println("Retrieved user_id from context:", userID)
	if !exists {
		c.Error(&errors.AppError{Message: "user_id not found in context", Code: http.StatusInternalServerError})
		return 0, false
	}
	id, ok := userID.(int64)
	fmt.Println("User ID type assertion success:", ok, "Value:", id)
	return id, ok
}
