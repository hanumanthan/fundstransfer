package service

import (
	"fundstransfer/database"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context){
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUser(c *gin.Context){
	var user models.User
	if err:= c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

