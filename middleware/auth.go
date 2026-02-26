package middleware

import (
	"REST-API/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authorization token required",
		})
		return
	}

	userID, role, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid or expired token",
		})
		return
	}

	context.Set("userId", userID)
	context.Set("role", role)

	context.Next()
}
