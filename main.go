package main

import (
	"fundstransfer/database"
	"fundstransfer/handlers"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()
	createTables()
	createUserAndAccount()
	registerRoutes(router)
	_ = router.Run()
}

func createTables() {
	database.DB.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Wallet{})

}

func registerRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Welcome to wallet funds transfer"})
	})

	router.GET("/liveness/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	router.GET("/users", handlers.GetUsers)
	router.GET("/wallets", handlers.GetWallets)
	router.POST("/transact", handlers.Transact)
	router.GET("/user/:user_id", handlers.GetUserDetails)
}

func createUserAndAccount() {
	for i := range []int{1, 2, 3} {
		user := &models.User{
			Name:     strconv.Itoa(i),
			Location: "singapore",
		}
		database.DB.Create(user)
		wallet := &models.Wallet{
			Balance:      100,
			UserId:       user.ID,
			MobileNumber: 11111111,
		}
		database.DB.Create(wallet)
	}
}
