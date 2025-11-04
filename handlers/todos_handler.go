package handlers

import (
	"auth-backend/database"
	"auth-backend/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Obtener todos los todos del usuario ordenados
func GetTodos(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdStr.(string))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.TodoCollection.Find(
		ctx,
		bson.M{"userId": userId},
		options.Find().SetSort(bson.M{"order": 1}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var todos []models.TodoItem
	if err := cursor.All(ctx, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Crear un todo nuevo
func CreateTodo(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdStr.(string))

	var todo models.TodoItem
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	todo.UserID = userId
	todo.Date = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := database.TodoCollection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todo.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, todo)
}

// Actualizar un todo existente
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	var todo models.TodoItem
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.TodoCollection.UpdateOne(
		ctx,
		bson.M{"_id": objId},
		bson.M{"$set": todo},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Borrar un todo
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.TodoCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Eliminado correctamente"})
}

// Reordenar todos los todos
func ReorderTodos(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdStr.(string))

	var newOrder []struct {
		ID    string `json:"id"`
		Order int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, item := range newOrder {
		objId, _ := primitive.ObjectIDFromHex(item.ID)
		_, err := database.TodoCollection.UpdateOne(
			ctx,
			bson.M{"_id": objId, "userId": userId},
			bson.M{"$set": bson.M{"order": item.Order}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reordenado correctamente"})
}
