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

// ✅ Obtener preferencias del usuario autenticado
func GetPreferencesHandler(c *gin.Context) {
	userID := c.MustGet("userId").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)

	collection := database.GetCollection("preferences")

	var prefs models.UserPreferences
	err := collection.FindOne(context.Background(), bson.M{"userId": objID}).Decode(&prefs)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"preferredLanguage": nil,
			"preferredTheme":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// ✅ Actualizar o crear preferencias del usuario
func UpdatePreferencesHandler(c *gin.Context) {
	userID := c.MustGet("userId").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)

	var body struct {
		PreferredLanguage string `json:"preferredLanguage"`
		PreferredTheme    string `json:"preferredTheme"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	collection := database.GetCollection("preferences")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🔹 Solo agregamos al update los campos enviados
	updateFields := bson.M{}
	if body.PreferredLanguage != "" {
		updateFields["preferredLanguage"] = body.PreferredLanguage
	}
	if body.PreferredTheme != "" {
		updateFields["preferredTheme"] = body.PreferredTheme
	}

	// Si no hay campos para actualizar, salimos
	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se enviaron campos válidos"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"userId": objID,
		},
	}
	for k, v := range updateFields {
		update["$set"].(bson.M)[k] = v
	}

	_, err := collection.UpdateOne(ctx, bson.M{"userId": objID}, update, options.Update().SetUpsert(true))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar preferencias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferencias actualizadas"})
}
