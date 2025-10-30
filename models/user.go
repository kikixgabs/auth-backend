package models

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Email    string `bson:"Email" json:"Email"`
	Password string `bson:"password" json:"password"`
}
