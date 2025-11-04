package main

import (
	"fmt"
	"os"

	"auth-backend/database"
	"auth-backend/handlers"
	"auth-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	mongoURI := os.Getenv("mongoURI")

	database.Connect(mongoURI, "todoappbd")

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas públicas
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)
	router.POST("/logout", handlers.LogoutHandler)

	// Rutas protegidas (ToDos + Preferencias)
	routes.RegisterRoutes(router)

	fmt.Println("🚀 Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}
