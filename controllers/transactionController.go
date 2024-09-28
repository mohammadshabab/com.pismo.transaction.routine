package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"com.pismo.transaction.routine/internal/database"
	"com.pismo.transaction.routine/models"
	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var newTransaction models.Transaction
	if err := c.BindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Negative amounts for purchases and withdrawals for operation 1 and 3
	if newTransaction.OperationTypeID == 1 || newTransaction.OperationTypeID == 3 {
		if newTransaction.Amount > 0 {
			newTransaction.Amount = -newTransaction.Amount
		}
	}

	newTransaction.EventDate = time.Now()

	result := database.DB.Create(&newTransaction)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create transaction"})
		return
	}
	slog.Info("Created a trnsaction ", "result ", result)
	c.JSON(http.StatusCreated, newTransaction)
}
