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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var expiration time.Duration
	if creds.RememberMe {
		expiration = 30 * 24 * time.Hour // 30 días
	} else {
		expiration = 24 * time.Hour // 1 día
	}

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
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, _ := token.SignedString(middleware.GetSecret())

	middleware.SetAuthCookie(c, tokenString, expiration)
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

func CheckEmailHandler(c *gin.Context) {
	var data struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"inUse": false})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": data.Email}).Decode(&user)

	c.JSON(http.StatusOK, gin.H{"inUse": err == nil})
}

func AuthMeHandler(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 🔍 Obtener la cookie donde guardás tu token
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// 🔐 Parsear JWT usando tu secret cargado en memory
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return middleware.GetSecret(), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			return
		}

		// 🆔 Obtener userId tolerante a variaciones
		var userIDStr string
		if v, ok := claims["userId"].(string); ok && v != "" {
			userIDStr = v
		} else if v, ok := claims["userID"].(string); ok && v != "" {
			userIDStr = v
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userId inválido"})
			return
		}

		// 🧪 Validar formato de Mongo ObjectID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de userId inválido"})
			return
		}

		// 🧑 Buscar usuario en DB
		var user models.User
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// 🔥 Respuesta final para tu frontend
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"user": gin.H{
				"id":       user.ID.Hex(),
				"email":    user.Email,
				"username": user.Username,
				"theme":    user.Theme,
				"language": user.Language,
			},
		})
	}
}
