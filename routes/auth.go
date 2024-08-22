package routes

import (
	"github.com/andymyp/to-do-list-go-api/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
	base := router.Group("/api/auth")
	base.POST("/register", controllers.Register)
}
