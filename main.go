package main

import (
	"fmt"
	"log"
	"os"

	"auth-backend/database"
	"auth-backend/handlers"
	"auth-backend/middleware"
	"auth-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 🧭 Detectar entorno
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// 🔑 Cargar .env adecuado
	if env == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Println("⚠️ No se pudo cargar .env.development, usando variables del sistema")
		}
	} else {
		if err := godotenv.Load(".env.production"); err != nil {
			log.Println("⚠️ No se pudo cargar .env.production, usando variables del sistema")
		}
	}

	middleware.LoadSecret()

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_NAME")

	// 💾 Conectar a MongoDB
	database.Connect(mongoURI, dbName)

	router := gin.Default()

	// 🌍 CORS dinámico según entorno
	allowedOrigins := []string{
		"http://localhost:4200",
		"https://kikixgabs.github.io",
	}

	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		for _, o := range allowedOrigins {
			if o == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", o)
				break
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 🧩 Rutas públicas
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)
	router.POST("/logout", handlers.LogoutHandler)
	router.POST("/check-email", handlers.CheckEmailHandler)
	router.GET("/auth/me", handlers.AuthMeHandler(database.UserCollection))

	// 🛣️ Rutas privadas (registradas en routes.go)
	routes.RegisterRoutes(router)

	// 🚀 Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("INFO: PORT not set, defaulting to " + port)
	}

	fmt.Printf("🚀 Servidor corriendo en modo %s en http://localhost:%s\n", env, port)
	router.Run(":" + port)
}
