package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"gochat/internal/config"
	"gochat/internal/database"
	"gochat/internal/handlers"
)

func main() {
    config.LoadEnv()
	database.Connect()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	port := config.GetPort()

	log.Println("Server running on port", port)
	r.Run(":" + port)
}