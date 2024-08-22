package models

import (
	"time"

	"github.com/andymyp/to-do-list-go-api/configs"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name" validate:"required"`
	Email     string             `json:"email" validate:"required,email"`
	Password  string             `json:"password" validate:"required,min=6"`
	Token     string             `json:"token"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func UserCollection() *mongo.Collection {
	return configs.DB.Collection("users")
}

type InputRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type InputLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
	Token string             `json:"token"`
}

type UserResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}

type Claims struct {
	ID primitive.ObjectID `json:"id"`
	jwt.RegisteredClaims
}
