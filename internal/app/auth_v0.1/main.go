package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/steve-mir/diivix_backend/internal/app/auth_v0.1/routes"
	// "github.com/steve-mir/diivix_backend/internal/app/auth_v0.1/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.Auth(router)
	routes.User(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
