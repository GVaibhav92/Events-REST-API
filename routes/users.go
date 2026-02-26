package routes

import (
	"REST-API/models"
	"REST-API/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "could not parse request data",
		})
		return
	}

	if validationErrors := utils.ValidateStruct(user); validationErrors != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	err = user.Save()
	if err != nil {
		if err.Error() == "email already registered" {
			context.JSON(http.StatusConflict, gin.H{
				"message": "email already registered",
			})
			return
		}
		// ADD THIS LINE TO SEE THE ACTUAL ERROR
		fmt.Println("[DEBUG] Error creating user:", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create user",
			"error":   err.Error(), // ALSO ADD THIS
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "could not parse request data",
		})
		return
	}

	if validationErrors := utils.ValidateStruct(user); validationErrors != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credentials",
		})
		return
	}

	// Generate access token (JWT)
	accessToken, err := utils.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate access token",
		})
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate refresh token",
		})
		return
	}

	// Save refresh token to database
	err = user.SaveRefreshToken(refreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not save refresh token",
		})
		return
	}

	// Return both tokens
	context.JSON(http.StatusOK, gin.H{
		"message":       "login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// refreshToken handles POST /auth/refresh
func refreshToken(context *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "refresh token required",
		})
		return
	}

	// Validate the refresh token and get the associated user
	user, err := models.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate access token",
		})
		return
	}

	// Rotate the refresh token
	newRefreshToken, err := user.RotateRefreshToken(request.RefreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not rotate refresh token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":       "token refreshed successfully",
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// logout handles POST /auth/logout
func logout(context *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "refresh token required",
		})
		return
	}

	// Delete the refresh token from database
	err = models.DeleteRefreshToken(request.RefreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not logout",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
