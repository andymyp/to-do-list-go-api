package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andymyp/to-do-list-go-api/configs"
	_ "github.com/andymyp/to-do-list-go-api/docs"
	"github.com/andymyp/to-do-list-go-api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//! For Generate Swagger Docs
// @title           							To-Do List API
// @version         							1.0
// @description     							Built with Go, Gin, MongoDB, JWT, and Swagger
// @contact.name   								API Support
// @contact.email  								andymyp1997@gmail.com
// @schemes 											http https
// @host      										localhost:3000
// @BasePath  										/api
// @securityDefinitions.apiKey  	Bearer
// @in 														header
// @name 													Authorization
// @description										Enter the token with the `Bearer prefix`, e.g. 'Bearer abcde12345'

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading Environments file!")
	}

	configs.ConnectDatabase()

	router := gin.Default()

	//! CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	//! Idle
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API server is running"})
	})

	//! API Docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//! All Routes
	routes.AuthRoute(router)
	routes.UserRoute(router)
	routes.TaskRoute(router)

	APP_PORT := os.Getenv("APP_PORT")
	APP_PORT = fmt.Sprintf(":%s", APP_PORT)

	router.Run(APP_PORT)
}
