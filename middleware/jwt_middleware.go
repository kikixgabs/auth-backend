package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 🔐 Clave secreta desde .env
var SecretKey = []byte(os.Getenv("JWT_SECRET"))

// Middleware de autenticación
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: token ausente"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userId", claims["userId"])
		c.Next()
	}
}

// Configura la cookie del token según entorno
func SetAuthCookie(c *gin.Context, tokenString string) {
	env := os.Getenv("APP_ENV")
	domain := "localhost"
	//secure := true

	if env == "production" {
		domain = "auth-backend-production-414c.up.railway.app"
		//secure = true
	}

	cookie := fmt.Sprintf(
		"token=%s; Path=/; Max-Age=%d; Domain=%s; HttpOnly; Secure; SameSite=None; Partitioned; Priority=High",
		tokenString,
		int(24*time.Hour.Seconds()),
		domain,
	)

	c.Writer.Header().Add("Set-Cookie", cookie)

	c.SetSameSite(http.SameSiteNoneMode)
}
