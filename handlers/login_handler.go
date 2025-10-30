package handlers

import (
	"auth-backend/database"
	"auth-backend/middleware"
	"auth-backend/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🔍 Buscar usuario por email
	var user models.User
	err := database.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// 🔑 Comparar contraseñas
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña incorrecta"})
		return
	}

	// 🪪 Crear token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(middleware.SecretKey)

	// 🍪 Guardar cookie segura
	middleware.SetAuthCookie(c, tokenString)

	c.JSON(http.StatusOK, gin.H{"message": "Login exitoso"})
}
