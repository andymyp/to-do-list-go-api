package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/andymyp/to-do-list-go-api/helpers"
	"github.com/andymyp/to-do-list-go-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(c *gin.Context) {
	var input models.InputTask
	ctx := context.Background()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	if err := helpers.ValidateStruct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)
	now := time.Now()

	createdAt, err := time.Parse("2006-01-02", input.CreatedAt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	deadlineAt, err := time.Parse("2006-01-02", input.DeadlineAt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	newTask := models.Task{
		UserID:      actualUser.ID,
		Title:       input.Title,
		Description: input.Description,
		CreatedAt:   createdAt,
		DeadlineAt:  deadlineAt,
		UpdatedAt:   now,
	}

	result, err := models.TaskCollection().InsertOne(ctx, newTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	var task models.Task
	taskID := result.InsertedID.(primitive.ObjectID)

	err = models.TaskCollection().FindOne(ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task created",
		"data":    task,
	})
}
