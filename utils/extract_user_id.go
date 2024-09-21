package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExtractUserID(c *gin.Context) (primitive.ObjectID, bool) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Status:  Error,
			Message: "Unauthorized",
		})
		return primitive.NilObjectID, false
	}

	userID, ok := userId.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  Error,
			Message: "Invalid user ID type",
		})
		return primitive.NilObjectID, false
	}

	return userID, true
}
