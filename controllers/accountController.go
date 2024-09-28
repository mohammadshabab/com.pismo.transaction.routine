package controllers

import (
	"net/http"
	"strconv"

	"com.pismo.transaction.routine/internal/database"
	"com.pismo.transaction.routine/models"
	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var newAccount models.Account
	if err := c.BindJSON(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&newAccount)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create account"})
		return
	}

	c.JSON(http.StatusCreated, newAccount)
}

func GetAccount(c *gin.Context) {
	accountId, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var account models.Account
	if result := database.DB.First(&account, accountId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}
