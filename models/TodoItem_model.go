package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Content   string             `bson:"content" json:"content"`
	Completed bool               `bson:"completed" json:"completed"`
	Priority  string             `bson:"priority" json:"priority"`
	Order     int                `bson:"order" json:"order"`
	Date      time.Time          `bson:"date" json:"date"`
}
