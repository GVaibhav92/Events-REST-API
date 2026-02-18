package main

import (
	"REST-API/config"
	"REST-API/db"
	"REST-API/middleware"
	"REST-API/routes"
	"REST-API/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load() //load the environment

	db.InitDB()
	utils.RegisterCustomValidations()

	server := gin.Default() //create Gin engine with Logger & Recovery middleware

	server.Use(middleware.Logger())
	routes.RegisterRoutes(server)

	server.Run(":" + config.App.Port)
}
