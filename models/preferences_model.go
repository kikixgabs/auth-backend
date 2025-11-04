package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPreferences struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID `bson:"userId" json:"userId"`
	PreferredLanguage string             `bson:"preferredLanguage,omitempty" json:"preferredLanguage,omitempty"`
	PreferredTheme    string             `bson:"preferredTheme,omitempty" json:"preferredTheme,omitempty"`
}
