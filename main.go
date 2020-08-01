package main

import (
	"fundstransfer/database"
	"fundstransfer/models"
	"fundstransfer/service"
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
	database.DB.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Account{})
	//database.DB.Model(&models.Account{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

}

func registerRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	router.GET("/liveness/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	router.GET("/users", service.GetUsers)
	router.GET("/accounts", service.GetAccounts)
	router.POST("/transact", createTransaction)
	router.GET("/user/:user_id", getUserDetails)
}

func createTransaction(c *gin.Context) {
	var createTransaction service.CreateTransaction
	if err := c.ShouldBindJSON(&createTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := createTransaction.CreateTransaction(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Transaction processed")
}

func getUserDetails(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("user_id"))
	userDetails, err := service.GetUserDetails(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": userDetails})
}

func createUserAndAccount() {
	for i:= range []int{1,2,3} {
		user := &models.User{
			Name:     strconv.Itoa(i),
			Location: "singapore",
		}
		database.DB.Create(user)
		account := &models.Account{
			Balance: 100,
			UserId:    user.ID,
		}
		database.DB.Create(account)
	}
}


