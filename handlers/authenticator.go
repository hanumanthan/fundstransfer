package handlers

import (
	"fmt"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Authenticate(apiType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("api_key")
		switch apiType {
		case "user":
			userId, _ := strconv.Atoi(c.Param("user_id"))
			if err := validateApiKeyForUser(userId, apiKey); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		case "admin":
			if err := validateApiKeyForAdmin(apiKey); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func validateApiKeyForUser(userId int, apiKey string) error {
	var credentials models.Credentials
	_ = credentials.GetApiKeyForUser(userId)
	if apiKey != credentials.ApiKey {
		return fmt.Errorf("invalid api key")
	}
	return nil
}

func validateApiKeyForAdmin(apiKey string) error {
	var user models.User
	var credentials models.Credentials
	_ = user.GetByName("admin")
	_ = credentials.GetApiKeyForUser(int(user.ID))
	if apiKey != credentials.ApiKey {
		return fmt.Errorf("invalid api key")
	}
	return nil
}
