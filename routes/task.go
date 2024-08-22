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
	base.GET("/tasks", controllers.GetTasks)
	base.GET("/task/:id", controllers.GetTask)
	base.PUT("/task/:id", controllers.UpdateTask)
	base.PUT("/task/status/:id", controllers.UpdateStatusTask)
	base.DELETE("/task/:id", controllers.DeleteTask)
}
