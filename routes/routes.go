package routes

import (
	"auth-backend/handlers"
	"auth-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Grupo de rutas protegidas (requieren autenticación)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Preferencias del usuario
		auth.GET("/preferences", handlers.GetPreferencesHandler)
		auth.PUT("/preferences", handlers.UpdatePreferencesHandler)

		// ToDos
		todos := auth.Group("/todos")
		{
			todos.GET("", handlers.GetTodos)
			todos.POST("", handlers.CreateTodo)
			todos.PUT("/:id", handlers.UpdateTodoHandler)
			todos.DELETE("/:id", handlers.DeleteTodo)
			todos.PUT("/reorder", handlers.ReorderTodos)
		}
	}
}
