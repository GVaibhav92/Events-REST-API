package routes

import (
	"REST-API/models"
	"REST-API/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	pageStr := context.DefaultQuery("page", "1")
	limitStr := context.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid page number",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid limit, must be between 1 and 100",
		})
		return
	}

	events, total, err := models.GetAllEvents(page, limit)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch events!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":       events,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit, // ceiling division
	})
}

func getEvent(context *gin.Context) {
	eventID := context.Param("id")
	id, err := strconv.Atoi(eventID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch event!",
		})
		return
	}

	if event == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "could not parse request data",
		})
		return
	}

	// Validate after binding
	if validationErrors := utils.ValidateStruct(event); validationErrors != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	userID, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not retrieve user identity",
		})
		return
	}
	event.UserID = userID.(int)

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create event!",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created",
		"event":   event,
	})
}

func updateEvent(context *gin.Context) {
	eventID := context.Param("id")
	id, err := strconv.Atoi(eventID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID",
		})
		return
	}

	existingEvent, err := models.GetEventByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch event!",
		})
		return
	}
	if existingEvent == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		return
	}

	userID, _ := context.Get("userId")
	if existingEvent.UserID != userID.(int) {
		context.JSON(http.StatusForbidden, gin.H{
			"message": "you are not authorized to update this event",
		})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "could not parse request data",
		})
		return
	}

	// Validate after binding
	if validationErrors := utils.ValidateStruct(updatedEvent); validationErrors != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	updatedEvent.ID = id
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not update event!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"event":   updatedEvent,
	})
}

func deleteEvent(context *gin.Context) {
	eventID := context.Param("id")
	id, err := strconv.Atoi(eventID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID",
		})
		return
	}

	existingEvent, err := models.GetEventByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch event!",
		})
		return
	}
	if existingEvent == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		return
	}

	// the logged in user owns this event
	userID, _ := context.Get("userId")
	if existingEvent.UserID != userID.(int) {
		context.JSON(http.StatusForbidden, gin.H{
			"message": "you are not authorized to delete this event",
		})
		return
	}

	err = existingEvent.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not delete event!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}
