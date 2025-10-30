package main

import (
	"fmt"

	"auth-backend/database"
	"auth-backend/handlers"
	"auth-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 🧠 Primero conectamos la base de datos
	database.Connect()

	router := gin.Default()

	// Rutas públicas
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.Login)

	// Rutas protegidas
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("", handlers.ProtectedHandler)

	fmt.Println("🚀 Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}
