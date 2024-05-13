package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required, min=2, max=100"`
	Last_name     *string            `json:"last_name" validate:"required, min=2, max=100"`
	Email         *string            `json:"email" validate:"required, email"`
	Password      *string            `json:"password" validate:"required, min=6, max=100"`
	Phone         *string            `json:"phone" validate:"required, min=10, max=14"`
	Token         *string            `json:"token"`
	Role          *string            `json:"role" validate:"required, eq=ADMIN|eq=USER|eq=ORG"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
