package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andymyp/to-do-list-go-api/configs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	APP_PORT := os.Getenv("APP_PORT")
	APP_PORT = fmt.Sprintf(":%s", APP_PORT)

	router.Run(APP_PORT)
}
