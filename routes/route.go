package routes

import (
	"REST-API/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// PUBLIC ROUTES (no auth required)
	server.POST("/signup", signup)
	server.POST("/login", login)
	server.POST("/auth/refresh", refreshToken)
	server.POST("/auth/logout", logout)

	// SEMI-PUBLIC ROUTES (anyone can view)
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	// PROTECTED ROUTES (authenticated users only)
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authenticate)
	{
		// Any logged-in user can create events
		authenticated.POST("/events", createEvent)

		// Only owner or admin can update/delete (checked in handler)
		authenticated.PUT("/events/:id", updateEvent)
		authenticated.DELETE("/events/:id", deleteEvent)

		// Any logged-in user can register for events
		authenticated.POST("/events/:id/register", registerForEvent)
		authenticated.DELETE("/events/:id/register", cancelRegistration)
	}

	// ADMIN-ONLY ROUTES (for future admin features)
	admin := server.Group("/admin")
	admin.Use(middleware.Authenticate, middleware.RequireAdmin)
	{
		// Future: admin.GET("/users", getAllUsers)
		// Future: admin.DELETE("/users/:id", deleteUser)
	}
}
