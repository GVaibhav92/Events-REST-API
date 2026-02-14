package routes

import (
	"REST-API/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userID, _ := context.Get("userId")

	eventID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid event ID",
		})
		return
	}

	registration := models.Registration{
		EventID: eventID,
		UserID:  userID.(int),
	}

	err = registration.Save()
	if err != nil {
		if err.Error() == "event not found" {
			context.JSON(http.StatusNotFound, gin.H{
				"message": "event not found",
			})
			return
		}
		if err.Error() == "already registered for this event" {
			context.JSON(http.StatusConflict, gin.H{
				"message": "you are already registered for this event",
			})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not register for event",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "successfully registered for event",
		"registration": gin.H{
			"id":      registration.ID,
			"eventId": registration.EventID,
			"userId":  registration.UserID,
		},
	})
}

func cancelRegistration(context *gin.Context) {
	userID, _ := context.Get("userId")

	eventID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid event ID",
		})
		return
	}

	registration := models.Registration{
		EventID: eventID,
		UserID:  userID.(int),
	}

	err = registration.Cancel()
	if err != nil {
		if err.Error() == "event not found" {
			context.JSON(http.StatusNotFound, gin.H{
				"message": "event not found",
			})
			return
		}
		if err.Error() == "you are not registered for this event" {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "you are not registered for this event",
			})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not cancel registration",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "registration cancelled successfully",
	})
}
