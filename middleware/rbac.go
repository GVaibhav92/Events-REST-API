package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// blocks requests if the user is not an admin
func RequireAdmin(context *gin.Context) {
	role, exists := context.Get("role")

	if !exists {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authentication required",
		})
		return
	}

	if role != "admin" {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "admin access required",
		})
		return
	}

	context.Next()
}

// checks if the user owns the resource OR is an admin
// resource owner ID should be set by the handler before this
func RequireOwnerOrAdmin(context *gin.Context) {
	role, roleExists := context.Get("role")
	userID, userExists := context.Get("userId")
	resourceOwnerID, ownerExists := context.Get("resourceOwnerId")

	if !roleExists || !userExists {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authentication required",
		})
		return
	}

	if role == "admin" {
		context.Next()
		return
	}

	if !ownerExists {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "resource owner not set",
		})
		return
	}

	if userID != resourceOwnerID {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "you don't have permission to access this resource",
		})
		return
	}

	context.Next()
}
