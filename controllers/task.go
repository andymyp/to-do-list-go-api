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

// CreateTask 		godoc
// @Security 			Bearer
// @Summary      	Create Task
// @Tags         	Task
// @Accept       	json
// @Produce      	json
// @Param        	request body models.InputTask true "Payload [Raw]"
// @Success      	200 "ok"
// @Router       	/task [post]
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

	var task models.TaskResponse
	taskID := result.InsertedID.(primitive.ObjectID)

	err = models.TaskCollection().FindOne(ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	task.StatusString = helpers.StatusString(task.Status)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task created",
		"data":    task,
	})
}

// GetTasks 			godoc
// @Security 			Bearer
// @Summary      	Get Tasks
// @Tags         	Task
// @Accept       	json
// @Produce      	json
// @Success      	200 "ok"
// @Router       	/tasks [get]
func GetTasks(c *gin.Context) {
	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	var tasks []models.TaskResponse
	ctx := context.Background()

	filter := bson.M{
		"user_id": actualUser.ID,
		"status": bson.M{
			"$ne": 3,
		},
	}

	cursor, err := models.TaskCollection().Find(ctx, filter)
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

		task.StatusString = helpers.StatusString(task.Status)

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

// GetTask 				godoc
// @Security 			Bearer
// @Summary      	Get Task
// @Tags         	Task
// @Accept       	json
// @Produce      	json
// @Param        	id path string true "Task ID"
// @Success      	200 "ok"
// @Router       	/task/{id} [get]
func GetTask(c *gin.Context) {
	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	var task models.TaskResponse
	ctx := context.Background()

	filter := bson.M{
		"_id":     id,
		"user_id": actualUser.ID,
		"status": bson.M{
			"$ne": 3,
		},
	}

	err = models.TaskCollection().FindOne(ctx, filter).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	task.StatusString = helpers.StatusString(task.Status)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "My Task",
		"data":    task,
	})
}

// UpdateTask 		godoc
// @Security 			Bearer
// @Summary      	Update Task
// @Tags         	Task
// @Accept       	json
// @Produce      	json
// @Param        	id path string true "Task ID"
// @Param        	request body models.InputTask true "Payload [Raw]"
// @Success      	200 "ok"
// @Router       	/task/{id} [put]
func UpdateTask(c *gin.Context) {
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

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

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

	filter := bson.M{
		"_id":     id,
		"user_id": actualUser.ID,
		"status": bson.M{
			"$ne": 3,
		},
	}

	updateTask := bson.M{
		"$set": bson.M{
			"title":       input.Title,
			"description": input.Description,
			"created_at":  createdAt,
			"deadline_at": deadlineAt,
			"updated_at":  time.Now(),
		},
	}

	_, err = models.TaskCollection().UpdateOne(ctx, filter, updateTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	var task models.TaskResponse

	err = models.TaskCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	task.StatusString = helpers.StatusString(task.Status)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task Updated",
		"data":    task,
	})
}

// UpdateStatusTask 		godoc
// @Security 						Bearer
// @Summary      				Update Status Task
// @Tags         				Task
// @Accept       				json
// @Produce      				json
// @Param        				id path string true "Task ID"
// @Param        				request body models.InputStatusTask true "Payload [Raw]"
// @Success      				200 "ok"
// @Router       				/task/status/{id} [put]
func UpdateStatusTask(c *gin.Context) {
	var input models.InputStatusTask
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

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	filter := bson.M{
		"_id":     id,
		"user_id": actualUser.ID,
		"status": bson.M{
			"$ne": 3,
		},
	}

	updateTask := bson.M{
		"$set": bson.M{
			"status":     input.Status,
			"updated_at": time.Now(),
		},
	}

	_, err = models.TaskCollection().UpdateOne(ctx, filter, updateTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	var task models.TaskResponse

	err = models.TaskCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	task.StatusString = helpers.StatusString(task.Status)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task Status Updated",
		"data":    task,
	})
}

// DeleteTask 		godoc
// @Security 			Bearer
// @Summary      	Delete Task
// @Tags         	Task
// @Accept       	json
// @Produce      	json
// @Param        	id path string true "Task ID"
// @Success      	200 "ok"
// @Router       	/task/{id} [delete]
func DeleteTask(c *gin.Context) {
	ctx := context.Background()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	filter := bson.M{
		"_id":     id,
		"user_id": actualUser.ID,
		"status": bson.M{
			"$ne": 3,
		},
	}

	updateTask := bson.M{
		"$set": bson.M{
			"status":     3,
			"updated_at": time.Now(),
		},
	}

	_, err = models.TaskCollection().UpdateOne(ctx, filter, updateTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	var task models.TaskResponse

	err = models.TaskCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	task.StatusString = helpers.StatusString(task.Status)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task Deleted",
		"data":    task,
	})
}
