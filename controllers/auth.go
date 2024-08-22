package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/andymyp/to-do-list-go-api/helpers"
	"github.com/andymyp/to-do-list-go-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	if err := helpers.ValidateStruct(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	var existingUser models.User
	ctx := context.Background()

	err := models.UserCollection().FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Email already in use, please login",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	now := time.Now()
	user.Password = hashedPassword
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err = models.UserCollection().InsertOne(ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var input models.InputLogin

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

	var user models.User
	ctx := context.Background()

	err := models.UserCollection().FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Email or password is incorrect",
		})
		return
	}

	if !helpers.CheckPassword(user.Password, input.Password) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Email or password is incorrect",
		})
		return
	}

	token, err := helpers.GenerateToken(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"token":      token,
			"updated_at": time.Now(),
		},
	}

	_, err = models.UserCollection().UpdateByID(ctx, user.ID, update)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	response := models.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: user.Token,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Login successful",
		"data":    response,
	})
}
