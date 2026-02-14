package routes

import (
	"REST-API/models"
	"REST-API/utils"
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
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create user",
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

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}
