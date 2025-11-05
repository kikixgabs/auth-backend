package main

import (
	"fmt"
	"log"
	"os"

	"auth-backend/database"
	"auth-backend/handlers"
	"auth-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// 🧭 Cargar variables según entorno
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Carga el .env correspondiente (solo en local)
	if env == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Println("⚠️ No se pudo cargar .env.development, usando variables del sistema")
		}
	} else {
		if err := godotenv.Load(".env.production"); err != nil {
			log.Println("⚠️ No se pudo cargar .env.production, usando variables del sistema")
		}
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_NAME")

	database.Connect(mongoURI, dbName)

	router := gin.Default()

	frontendURL := "http:/localhost:4200"
	if env == "production" {
		frontendURL = "https://kikixgabs.github.io"
	}

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("INFO: PORT not set, defaulting to " + port)
	}

	log.Println("Server starting on port " + port)

	routes.RegisterRoutes(router)

	fmt.Println("🚀 Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}
