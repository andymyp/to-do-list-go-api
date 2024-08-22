package models

import (
	"time"

	"github.com/andymyp/to-do-list-go-api/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description" validate:"required"`
	Status      int                `bson:"status" json:"status"` // 0: waiting list, 1: in progress, 2: done, 3: deleted
	CreatedAt   time.Time          `bson:"created_at" json:"created_at" validate:"required"`
	DeadlineAt  time.Time          `bson:"deadline_at" json:"deadline_at" validate:"required"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

func TaskCollection() *mongo.Collection {
	return configs.DB.Collection("tasks")
}

type InputTask struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	CreatedAt   string `bson:"created_at" json:"created_at" validate:"required"`
	DeadlineAt  string `bson:"deadline_at" json:"deadline_at" validate:"required"`
}

type InputStatusTask struct {
	Status int `json:"status" validate:"required"`
}

type TaskResponse struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Status       int                `json:"status"`
	StatusString string             `json:"status_string"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	DeadlineAt   time.Time          `bson:"deadline_at" json:"deadline_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
