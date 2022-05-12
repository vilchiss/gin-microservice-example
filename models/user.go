package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters User signin
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	RegisteredAt time.Time          `json:"registered_at" bson:"registered_at"`
	Password     string             `json:"password" bson:"password"`
	Username     string             `json:"username" bson:"username"`
}
