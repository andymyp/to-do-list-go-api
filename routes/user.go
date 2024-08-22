package routes

import (
	"github.com/andymyp/to-do-list-go-api/controllers"
	"github.com/andymyp/to-do-list-go-api/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	base := router.Group("/api/user")
	base.Use(middlewares.AuthMiddleware())
	base.GET("/profile", controllers.UserProfile)
}
