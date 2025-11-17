package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// secretKey se carga desde main.go llamando a LoadSecret()
var secretKey []byte

// LoadSecret se llama desde main.go después de cargar .env
func LoadSecret() {
	secretKey = []byte(os.Getenv("JWT_SECRET"))
}

// GetSecret devuelve la clave cargada
func GetSecret() []byte {
	return secretKey
}

// Middleware de autenticación
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: token ausente"})
			c.Abort()
			return
		}

		// Usamos la misma secret cargada globalmente
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return GetSecret(), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			c.Abort()
			return
		}

		// Intentar obtener userId con tolerancia a variaciones
		var userIDStr string
		if v, ok := claims["userId"].(string); ok && v != "" {
			userIDStr = v
		} else if v, ok := claims["userID"].(string); ok && v != "" {
			userIDStr = v
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userId inválido"})
			c.Abort()
			return
		}

		// validar formato de ObjectID
		if _, err := primitive.ObjectIDFromHex(userIDStr); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de userId inválido"})
			c.Abort()
			return
		}

		// Guardar claim en el contexto usando la clave canonical "userId"
		c.Set("userId", userIDStr)
		c.Next()
	}
}

// SetAuthCookie ahora usa c.SetCookie para manejar flags correctamente
func SetAuthCookie(c *gin.Context, tokenString string, duration time.Duration) {
	env := os.Getenv("APP_ENV")

	maxAge := int(duration.Seconds())
	domain := ""
	secure := false

	// En producción configurar dominio y secure
	if env == "production" {
		domain = "auth-backend-production-414c.up.railway.app" // adaptá si tenés otro host
		secure = true
	}

	// name, value, maxAge, path, domain, secure, httpOnly
	c.SetCookie("token", tokenString, maxAge, "/", domain, secure, true)

	// SameSite None para permitir cross-site if needed (ya se configura así)
	c.SetSameSite(http.SameSiteNoneMode)
}
