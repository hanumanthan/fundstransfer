package main

import (
	"fundstransfer/handlers"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	path := "payments.db"
	_ = os.Remove(path)
	models.ConnectDatabase()
	models.CreateTables()
	models.CreateUserAndWallet()
	registerRoutes(router)
	_ = router.Run()
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

	}


}
