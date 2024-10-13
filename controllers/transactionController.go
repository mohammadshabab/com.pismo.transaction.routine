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

	// Binding the JSON request body to the Transaction model
	if err := c.BindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if account id exist before proceeding with transaction creation
	var account models.Account
	if err := database.Repo.RDB.First(&account, newTransaction.AccountID).Error; err != nil {
		slog.Error("account not found ", "account_id ", newTransaction.AccountID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "account does not exist"})
		return
	}

	// Starting a new database transaction
	result := database.Repo.RDB.Begin()
	if err := result.Error; err != nil {
		slog.Error("Failed to start a transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to start transaction"})
		return
	}

	// Handle error and ensure transaction rollback
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic occurred while creating a transaction", "panic", r)
			result.Rollback()
		}
	}()

	// Negative amounts for purchases and withdrawals (operation type 1 and 3)
	if newTransaction.OperationTypeID == 1 || newTransaction.OperationTypeID == 3 {
		if newTransaction.Amount > 0 {
			newTransaction.Amount = -newTransaction.Amount
		}
	}

	// Set the event date to the current time
	newTransaction.EventDate = time.Now()

	// Create the transaction in the database
	if err := result.Create(&newTransaction).Error; err != nil {
		slog.Error("Failed to create transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create transaction"})
		result.Rollback()
		return
	}

	// Commit the transaction after successful creation
	if err := result.Commit().Error; err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to commit transaction"})
		return
	}

	slog.Info("Transaction created successfully", "transaction_id", newTransaction.ID, "operation_type", newTransaction.OperationTypeID)
	c.JSON(http.StatusCreated, newTransaction)
}
