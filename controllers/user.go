package controllers

import (
	"context"
	"net/http"

	"github.com/andymyp/to-do-list-go-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// UserProfile 		godoc
// @Security 			Bearer
// @Summary      	User profile
// @Tags         	User
// @Accept       	json
// @Produce      	json
// @Success      	200 "ok"
// @Router       	/user/profile [get]
func UserProfile(c *gin.Context) {
	getuser, _ := c.Get("user")
	actualUser, _ := getuser.(models.UserResponse)

	var user models.User
	ctx := context.Background()

	err := models.UserCollection().FindOne(ctx, bson.M{"_id": actualUser.ID}).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
		return
	}

	response := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   response,
	})
}
