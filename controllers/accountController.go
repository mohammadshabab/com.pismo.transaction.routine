package controllers

import (
	"log/slog"
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
		slog.Error("Bad request ", "error: ", err)
		return
	}

	result := database.Repo.RDB.Create(&newAccount)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create account"})
		slog.Error("Unable to create account ", "error: ", result.Error)
		return
	}
	slog.Info("New account has been created ", "account info: ", newAccount)
	c.JSON(http.StatusCreated, newAccount)
}

func GetAccount(c *gin.Context) {
	accountId, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		slog.Error("Account id must be an integer ", "accountId ", accountId, "err ", err)
		return
	}

	var account models.Account
	if result := database.Repo.RDB.First(&account, accountId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		slog.Error("account_id does not exist ", "account_id: ", accountId, "err ", err)
		return
	}
	slog.Info("account detail ", "acc: ", account)
	c.JSON(http.StatusOK, account)
}
