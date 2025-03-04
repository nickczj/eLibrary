package routes

import (
	"eLibrary/global"
	"eLibrary/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock DB setup
func setupMockDB() {
	mockDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	global.Database = mockDB

	// Migrate the schema
	global.Database.AutoMigrate(&model.BookDetail{})
	global.Database.AutoMigrate(&model.User{})
	global.Database.AutoMigrate(&model.LoanDetail{})

	// Seed test data
	err := global.Database.Create(&model.BookDetail{
		Model: gorm.Model{
			ID: 1,
		},
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "1234567890",
		AvailableCopies: 5,
	}).Error
	if err != nil {
		panic("failed to create test book")
	}

	err = global.Database.Create(&model.BookDetail{
		Title:           "Second Book",
		Author:          "Second Author",
		ISBN:            "1234567899",
		AvailableCopies: 0,
	}).Error
	if err != nil {
		panic("failed to create second test book")
	}

	err = global.Database.Create(&model.User{
		Model: gorm.Model{
			ID: 1,
		},
		FirstName: "Nick",
		LastName:  "Chow",
		Username:  "nickczj",
		Email:     "nick.chow.zj@gmail.com",
	}).Error
	if err != nil {
		panic("failed to create user")
	}
}

func TestGetBookAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setupMockDB()

	router := SetupRouter()

	t.Run("Book Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/elibrary/v1/book/Test%20Book", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]model.BookDetail
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Book", response["book"].Title)
	})

	t.Run("Book Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/elibrary/v1/book/Unknown%20Book", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)

		// Validate that there's no "book" field or that it is empty
		assert.Nil(t, response["book"])
	})
}

func TestBorrowBookAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock database or mock global.Database
	setupMockDB()

	router := SetupRouter()

	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/elibrary/v1/borrow", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "invalid request body")
	})

	t.Run("Validation Failed", func(t *testing.T) {
		invalidLoan := `{
            "title": "",
            "user_id": 0,
            "name_of_borrower": ""
        }`

		req, _ := http.NewRequest("POST", "/elibrary/v1/borrow", strings.NewReader(invalidLoan))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "validation failed")
	})

	t.Run("Book Not Available", func(t *testing.T) {
		// Prepare loan request for a book that has no available copies
		reqBody := `{
            "title": "Second Book",
            "user_id": 1,
            "name_of_borrower": "John Doe"
        }`

		// Simulate a scenario where no copies are available for the book
		req, _ := http.NewRequest("POST", "/elibrary/v1/borrow", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "there are no more available books to borrow")
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Prepare loan request for a book and a non-existent user
		reqBody := `{
            "title": "Test Book",
            "user_id": 99999
        }`

		// Simulate a scenario where the user doesn't exist
		req, _ := http.NewRequest("POST", "/elibrary/v1/borrow", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "unable to borrow book")
	})

	t.Run("Successful Borrow", func(t *testing.T) {
		reqBody := `{
            "title": "Test Book",
            "user_id": 1
        }`

		// Simulate a successful borrowing scenario where a book is available and user exists
		req, _ := http.NewRequest("POST", "/elibrary/v1/borrow", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		fmt.Println("Response Status:", resp.Code)                              // Print the response code
		fmt.Println("Response Location Header:", resp.Header().Get("Location")) // Print the redirect URL

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that the loan response contains expected data
		loan := response["loan"].(map[string]interface{})
		assert.NotNil(t, loan)
		assert.Equal(t, "Nick Chow", loan["name_of_borrower"])
	})
}

func TestExtendBookAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock database or mock global.Database
	setupMockDB()

	router := SetupRouter()

	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/elibrary/v1/extend", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "invalid request body")
	})

	t.Run("Validation Failed", func(t *testing.T) {
		invalidLoan := `{
            "title": "",
            "user_id": 0,
            "name_of_borrower": ""
        }`

		req, _ := http.NewRequest("POST", "/elibrary/v1/extend", strings.NewReader(invalidLoan))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "validation failed")
	})

	t.Run("Loan Not Found", func(t *testing.T) {
		// Prepare loan request for a book and a non-existent user
		reqBody := `{
            "title": "Test Book",
            "user_id": 99999
        }`

		// Simulate a scenario where the user doesn't exist
		req, _ := http.NewRequest("POST", "/elibrary/v1/extend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "loan not found")
	})

	t.Run("Successful Extend", func(t *testing.T) {
		err := global.Database.Create(&model.LoanDetail{
			BookID: 1,
			BookDetail: model.BookDetail{
				Model: gorm.Model{
					ID: 1,
				},
				Title:           "Test Book",
				AvailableCopies: 5,
			},
			UserID: 1,
			User: model.User{
				Model: gorm.Model{
					ID: 1,
				},
				FirstName: "Nick",
				LastName:  "Chow",
			},
			NameOfBorrower: "Nick Chow",
			LoanDate:       time.Now(),
			ReturnDate:     time.Now().AddDate(0, 0, 28),
			IsReturned:     false,
		}).Error
		if err != nil {
			panic("failed to create second test book")
		}

		reqBody := `{
            "title": "Test Book",
            "user_id": 1
        }`

		// Simulate a successful borrowing scenario where a book is available and user exists
		req, _ := http.NewRequest("POST", "/elibrary/v1/extend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that the loan response contains expected data
		loan := response["loan"].(map[string]interface{})
		assert.NotNil(t, loan)
		assert.Equal(t, "Nick Chow", loan["name_of_borrower"])
	})
}

func TestReturnBookAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock database or mock global.Database
	setupMockDB()

	router := SetupRouter()

	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/elibrary/v1/return", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "invalid request body")
	})

	t.Run("Validation Failed", func(t *testing.T) {
		invalidLoan := `{
            "title": "",
            "user_id": 0,
            "name_of_borrower": ""
        }`

		req, _ := http.NewRequest("POST", "/elibrary/v1/return", strings.NewReader(invalidLoan))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "validation failed")
	})

	t.Run("Loan Not Found", func(t *testing.T) {
		// Prepare loan request for a book and a non-existent user
		reqBody := `{
            "title": "Test Book",
            "user_id": 99999
        }`

		// Simulate a scenario where the user doesn't exist
		req, _ := http.NewRequest("POST", "/elibrary/v1/return", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["bad request"], "loan not found")
	})

	t.Run("Successful Return", func(t *testing.T) {
		err := global.Database.Create(&model.LoanDetail{
			BookID: 1,
			BookDetail: model.BookDetail{
				Model: gorm.Model{
					ID: 1,
				},
				Title:           "Test Book",
				AvailableCopies: 5,
			},
			UserID: 1,
			User: model.User{
				Model: gorm.Model{
					ID: 1,
				},
				FirstName: "Nick",
				LastName:  "Chow",
			},
			NameOfBorrower: "Nick Chow",
			LoanDate:       time.Now(),
			ReturnDate:     time.Now().AddDate(0, 0, 28),
			IsReturned:     false,
		}).Error
		if err != nil {
			panic("failed to create second test book")
		}

		reqBody := `{
            "title": "Test Book",
            "user_id": 1
        }`

		// Simulate a successful borrowing scenario where a book is available and user exists
		req, _ := http.NewRequest("POST", "/elibrary/v1/return", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that the loan response contains expected data
		loan := response["loan"].(map[string]interface{})
		assert.NotNil(t, loan)
		assert.Equal(t, "Nick Chow", loan["name_of_borrower"])
	})
}
