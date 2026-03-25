package main

import (
	"REST-API/config"
	"REST-API/db"
	"REST-API/middleware"
	"REST-API/routes"
	"REST-API/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	db.InitDB()
	utils.RegisterCustomValidations()

	server := gin.Default()

	server.Use(middleware.RequestID)
	server.Use(middleware.Timeout(config.App.RequestTimeout))
	server.Use(middleware.Logger())

	routes.RegisterRoutes(server)

	httpServer := &http.Server{
		Addr:    ":" + config.App.Port,
		Handler: server,
	}

	// goroutine to not block signal handling
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	//block until SIGINT or SIGTERM is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	//allow up to 10 seconds for active requests to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// close db connection after all requests have finished
	if err := db.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Shutdown complete")
}
