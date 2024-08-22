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

func MyTasks(c *gin.Context) {
	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	var tasks []models.TaskResponse
	ctx := context.Background()

	cursor, err := models.TaskCollection().Find(ctx, bson.M{"user_id": actualUser.ID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	for cursor.Next(ctx) {
		var task models.TaskResponse
		if err := cursor.Decode(&task); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		switch task.Status {
		case 0:
			task.StatusString = "Waiting List"
		case 1:
			task.StatusString = "In Progress"
		case 2:
			task.StatusString = "Done"
		default:
			task.StatusString = "Deleted"
		}

		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "My Tasks",
		"data":    tasks,
	})
}

func MyTask(c *gin.Context) {
	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var task models.TaskResponse
	ctx := context.Background()

	filter := bson.M{
		"_id":     id,
		"user_id": actualUser.ID,
	}

	err := models.TaskCollection().FindOne(ctx, filter).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	switch task.Status {
	case 0:
		task.StatusString = "Waiting List"
	case 1:
		task.StatusString = "In Progress"
	case 2:
		task.StatusString = "Done"
	default:
		task.StatusString = "Deleted"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "My Task",
		"data":    task,
	})
}
