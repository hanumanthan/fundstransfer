package main

import (
	"fundstransfer/database"
	"fundstransfer/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()
	registerRoutes(router)
	router.Run()
}

func registerRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	router.GET("/liveness/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	router.GET("/users", service.GetUsers)
	router.POST("/users", service.CreateUser)
}
