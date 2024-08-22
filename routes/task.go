package routes

import (
	"github.com/andymyp/to-do-list-go-api/controllers"
	"github.com/andymyp/to-do-list-go-api/middlewares"
	"github.com/gin-gonic/gin"
)

func TaskRoute(router *gin.Engine) {
	base := router.Group("/api")
	base.Use(middlewares.AuthMiddleware())
	base.POST("/task", controllers.CreateTask)
}
