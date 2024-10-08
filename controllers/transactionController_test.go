package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.pismo.transaction.routine/internal/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize sqlmock
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	defer mockDB.Close()

	sqlMock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("11.3.2"))

	// Open a gorm.DB connection using sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: mockDB, // Use sqlmock's database connection here
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm DB: %v", err)
	}

	// Create a repository with the mocked DB
	database.Repo = database.NewRepository(gormDB)

	// Test cases
	tests := []struct {
		name               string
		input              interface{}
		expectedStatusCode int
		expectedResponse   gin.H
		mockExpectations   func()
	}{
		{
			name:               "Success case for operation id 4",
			input:              gin.H{"account_id": 1, "operation_type_id": 4, "amount": 250.45},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"account_id": float64(1), "operation_type_id": float64(4), "amount": float64(250.45)},
			mockExpectations: func() {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("INSERT INTO `transactions`"). // Adjust the query according to your actual query structure
											WithArgs(int(1), int(4), float64(250.45), sqlmock.AnyArg()).
											WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
			},
		},
		{
			name:               "Success case for operation id 2",
			input:              gin.H{"account_id": 1, "operation_type_id": 2, "amount": 250.45},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"account_id": float64(1), "operation_type_id": float64(2), "amount": float64(250.45)},
			mockExpectations: func() {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("INSERT INTO `transactions`"). // Adjust the query according to your actual query structure
											WithArgs(1, 2, float64(250.45), sqlmock.AnyArg()).
											WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
			},
		},
		{
			name:               "Success case for operation id 1",
			input:              gin.H{"account_id": 1, "operation_type_id": 1, "amount": 250.45},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"account_id": float64(1), "operation_type_id": float64(1), "amount": -float64(250.45)},
			mockExpectations: func() {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("INSERT INTO `transactions`"). // Adjust the query according to your actual query structure
											WithArgs(1, 1, -float64(250.45), sqlmock.AnyArg()).
											WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
			},
		},
		{
			name:               "Success case for operation id 3",
			input:              gin.H{"account_id": 1, "operation_type_id": 3, "amount": 250.45},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"account_id": float64(1), "operation_type_id": float64(3), "amount": -float64(250.45)},
			mockExpectations: func() {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("INSERT INTO `transactions`"). // Adjust the query according to your actual query structure
											WithArgs(int(1), int(3), -float64(250.45), sqlmock.AnyArg()).
											WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
			},
		},
		{
			name:               "Invalid JSON binding",
			input:              "invalid-json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"error": "json: cannot unmarshal string into Go value of type models.Transaction"},
		},
		{
			name:               "Failed to start transaction",
			input:              gin.H{"account_id": 1, "operation_type_id": 4, "amount": float64(250.45)},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   gin.H{"error": "Unable to start transaction"},
			mockExpectations: func() {
				sqlMock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
			},
		},
		{
			name:               "Failed to create transaction",
			input:              gin.H{"account_id": 1, "operation_type_id": 4, "amount": float64(250.45)},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   gin.H{"error": "Unable to create transaction"},
			mockExpectations: func() {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("INSERT INTO `transactions`").
					WithArgs(1, 4, float64(250.45), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)
			},
		},
		{
			name:               "Failed to commit transaction",
			input:              gin.H{"account_id": 1, "operation_type_id": 4, "amount": float64(250.45)},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   gin.H{"error": "Unable to commit transaction"},
			mockExpectations: func() {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("INSERT INTO `transactions`").
					WithArgs(1, 4, float64(250.45), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit().WillReturnError(gorm.ErrInvalidTransaction)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}

			// Create the Gin context and recorder
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			// Prepare the request body
			jsonBody, _ := json.Marshal(tt.input)
			request, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
			request.Header.Set("Content-Type", "application/json")
			context.Request = request

			// Call the handler
			CreateTransaction(context)

			// Check the status code and response
			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Nil(t, err)

			// Assert response
			for key, value := range tt.expectedResponse {
				assert.Equal(t, value, response[key])
			}

			// Assert all expectations were met
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
