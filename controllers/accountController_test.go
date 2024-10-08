package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"com.pismo.transaction.routine/internal/database"
	"com.pismo.transaction.routine/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	gormDB  *gorm.DB
	sqlMock sqlmock.Sqlmock
	mockDB  *sql.DB
)

func setupTestDB(t *testing.T) {
	var err error
	mockDB, sqlMock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	sqlMock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("11.3.2"))

	// Open a gorm.DB connection using sqlmock
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm DB: %v", err)
	}

	// Replace the real repository with the mock
	database.Repo = database.NewRepository(gormDB)
}

func TestGetAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup the test DB with sqlmock
	setupTestDB(t)

	// Set up the router and route
	router := gin.Default()
	router.GET("/accounts/:accountId", GetAccount)

	// Test case 1: Valid account ID (Account found)
	t.Run("ValidAccountID", func(t *testing.T) {
		accountId := 1

		// Set up the expected SQL query (ensure it matches exactly what GORM generates)
		sqlMock.ExpectQuery("^SELECT \\* FROM `accounts` WHERE `accounts`.`account_id` = \\? ORDER BY `accounts`.`account_id` LIMIT \\?$").
			WithArgs(1, 1). // Provide both account_id and LIMIT as arguments
			WillReturnRows(sqlmock.NewRows([]string{"account_id", "document_number"}).
				AddRow(1, "123456789"))

		// Create the HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/accounts/"+strconv.Itoa(accountId), nil)
		resp := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(resp, req)

		// Check the response
		assert.Equal(t, http.StatusOK, resp.Code)
		expected := `{"account_id":1,"document_number":"123456789"}`
		assert.JSONEq(t, expected, resp.Body.String())

		// Check expectations after the test
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Logf("Expectations after ValidAccountID: %s", err)
		}
	})

	// Test case 2: Invalid account ID (Non-integer)
	t.Run("InvalidAccountID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/accounts/abc", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		expected := `{"error":"Invalid account ID"}`
		assert.JSONEq(t, expected, resp.Body.String())

		// Check expectations after the test
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Logf("Expectations after InvalidAccountID: %s", err)
		}
	})

	// Test case 3: Account not found
	t.Run("AccountNotFound", func(t *testing.T) {
		accountId := 999

		// Use the exact SQL string instead of a regex.
		sqlMock.ExpectQuery("SELECT \\* FROM `accounts` WHERE `accounts`.`account_id` = ? ORDER BY `accounts`.`account_id` LIMIT 1").
			WithArgs(accountId).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/accounts/"+strconv.Itoa(accountId), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		expected := `{"error":"Account not found"}`
		assert.JSONEq(t, expected, resp.Body.String())

		// Check expectations after the test
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Logf("Expectations after AccountNotFound: %s", err)
		}
	})

}

//==================

func TestCreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up the test DB with sqlmock
	setupTestDB(t)

	// Set up the router and route
	router := gin.Default()
	router.POST("/accounts", CreateAccount)

	// Test case 1: Valid input
	t.Run("ValidInput", func(t *testing.T) {
		newAccount := models.Account{
			AccountID:      1,
			DocumentNumber: "123456789",
		}
		jsonBody, _ := json.Marshal(newAccount)

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("INSERT INTO `accounts`").WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.JSONEq(t, string(jsonBody), resp.Body.String())
	})

	// Test case 2: Invalid input (malformed JSON)
	t.Run("InvalidInput", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer([]byte("invalid json")))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		expected := `{"error":"invalid character 'i' looking for beginning of value"}`
		assert.JSONEq(t, expected, resp.Body.String())
	})

	// Test case 3: Database error
	t.Run("DatabaseError", func(t *testing.T) {
		newAccount := models.Account{
			AccountID:      1,
			DocumentNumber: "123456789",
		}
		jsonBody, _ := json.Marshal(newAccount)

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("INSERT INTO `accounts`").WillReturnError(fmt.Errorf("some database error"))
		sqlMock.ExpectRollback()

		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		expected := `{"error":"Unable to create account"}`
		assert.JSONEq(t, expected, resp.Body.String())
	})

	// Verify that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
