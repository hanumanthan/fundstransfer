package handlers

import (
	"fundstransfer/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Get all wallets in the system
func GetWallets(c *gin.Context) {
	wallets, err := models.GetAllWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, &Wallets{Wallets:wallets})
}
