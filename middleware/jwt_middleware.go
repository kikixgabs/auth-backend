package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("clave_super_secreta")

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

func SetAuthCookie(c *gin.Context, tokenString string) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"token",
		tokenString,
		int(24*time.Hour.Seconds()),
		"/",
		"auth-backend-production-414c.up.railway.app", // dominio explícito
		true, // secure = true, porque estás en https
		true, // httpOnly
	)

}
