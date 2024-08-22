package middlewares

import (
	"context"
	"net/http"
	"os"

	"github.com/andymyp/to-do-list-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var secretKeyEnv = os.Getenv("JWT_SECRET_KEY")
var jwtSecretKey = []byte(secretKeyEnv)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Not authorized!",
			})
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid or expired token!",
			})
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid token claims!",
			})
			return
		}

		var user models.User
		ctx := context.Background()

		filter := bson.M{
			"_id":   claims.ID,
			"token": tokenString,
		}

		if err := models.UserCollection().FindOne(ctx, filter).Decode(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"filter":  filter,
				"err":     err.Error(),
				"message": "Invalid or expired token!",
			})
			return
		}

		response := models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}

		c.Set("user", response)
		c.Set("claims", claims)
		c.Next()
	}
}
