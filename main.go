package main

import (
	"fmt"
	"fundstransfer/pkg/handlers"
	"fundstransfer/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func init() {
	models.ConnectDatabase()
	models.CreateTables()
	models.SetupData()
}

func registerRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Welcome to wallet funds transfer"})
	})

	router.GET("/liveness/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	userAuth := router.Group("/api/v1")
	userAuth.Use(handlers.Authenticate("user"))
	{
		userAuth.POST("/user/:user_id/transact", handlers.Transact)
		userAuth.GET("/user/:user_id", handlers.GetUserDetails)

	}

	adminAuth := router.Group("/admin")
	adminAuth.Use(handlers.Authenticate("admin"))
	{
		adminAuth.GET("/users", handlers.GetUsers)
		adminAuth.GET("/wallets", handlers.GetWallets)
		adminAuth.GET("/transactions", handlers.GetTransactions)

	}
}

func main() {
	router := gin.Default()
	registerRoutes(router)
	if err := router.Run(); err != nil {
		fmt.Println("Error starting gin router")
		os.Exit(1)
	}
}
