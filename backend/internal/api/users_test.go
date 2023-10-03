package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alhazenlabs/procure-hub/backend/internal/database"
	"github.com/alhazenlabs/procure-hub/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	// Set up a test Gin router
	router := gin.Default()

	// Initialize the database
	database.InitDB("sqlite3", ":memory:")
	database.DB.AutoMigrate(&models.User{})

	// Define a test case
	testUser := models.User{
		Email:       "test@example.com",
		Password:    "testpassword",
		UserType:    models.PrimaryOwner,
		CompanyName: "Test Company",
	}

	t.Run("RegisterUser_Success", func(t *testing.T) {
		// Create a request body with the testUser data
		requestBody, err := json.Marshal(testUser)
		assert.NoError(t, err)

		// Create a POST request to the /users/v1/signup endpoint with the request body
		req, err := http.NewRequest("POST", "/users/v1/signup", bytes.NewReader(requestBody))
		assert.NoError(t, err)

		// Set the request content type to JSON
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder to capture the response
		w := httptest.NewRecorder()

		// Create a Gin context from the request and response recorder
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the RegisterUser function to handle the request
		RegisterUser(c)

		// Assert the response status code
		assert.Equal(t, http.StatusCreated, w.Code)

		// Assert the response body
		var responseBody map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", responseBody["message"])
	})

	t.Run("RegisterUser_ValidationError", func(t *testing.T) {
		// Create a request body with missing required fields
		invalidUser := models.User{
			Email: "invalid@example.com",
		}
		requestBody, err := json.Marshal(invalidUser)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/users/v1/signup", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		RegisterUser(c)

		// Assert the response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Add more test cases as needed
}

func TestLoginHandler(t *testing.T) {
	// Set up a test Gin router
	router := gin.Default()

	// Initialize the database
	database.InitDB("sqlite3", ":memory:")
	database.DB.AutoMigrate(&models.User{})

	// Create a test user for login
	testUser := models.User{
		Email:       "test@example.com",
		Password:    "testpassword",
		UserType:    models.PrimaryOwner,
		CompanyName: "Test Company",
	}
	testUser.PasswordHash, _ = hashPassword(testUser.Password)
	database.DB.Create(&testUser)

	t.Run("LoginHandler_Success", func(t *testing.T) {
		// Create a request body with the testUser's email and password
		loginData := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    testUser.Email,
			Password: "testpassword",
		}
		requestBody, err := json.Marshal(loginData)
		assert.NoError(t, err)

		// Create a POST request to the /users/v1/login endpoint with the request body
		req, err := http.NewRequest("POST", "/users/v1/login", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder to capture the response
		w := httptest.NewRecorder()

		// Create a Gin context from the request and response recorder
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the LoginHandler function to handle the request
		LoginHandler(c)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("LoginHandler_InvalidCredentials", func(t *testing.T) {
		// Create a request body with invalid credentials
		invalidLoginData := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    testUser.Email,
			Password: "incorrectpassword",
		}
		requestBody, err := json.Marshal(invalidLoginData)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/users/v1/login", bytes.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		LoginHandler(c)

		// Assert the response status code
		// Assert the response status code
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Assert the response body (optional)
		var responseBody map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid credentials", responseBody["error"])
	})
}
