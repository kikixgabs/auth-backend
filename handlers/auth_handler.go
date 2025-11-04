package handlers

import (
	"auth-backend/database"
	"auth-backend/models"
	"context"
	"net/http"
	"time"

	"auth-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verificar si el email ya existe
	var existing models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El email ya está registrado"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar contraseña"})
		return
	}
	user.Password = string(hash)

	_, err = database.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado con éxito"})
}

func LoginHandler(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Generar JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(middleware.SecretKey)
	middleware.SetAuthCookie(c, tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Logueado correctamente"})
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie(
		"token", // nombre
		"",      // valor vacío
		-1,      // expiración inmediata
		"/",     // path
		"",      // dominio (vacío = mismo dominio)
		false,   // secure
		true,    // httpOnly
	)
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada correctamente"})
}
