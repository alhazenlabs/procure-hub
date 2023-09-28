package api

import (
	"net/http"

	"github.com/alhazenlabs/procure-hub/backend/internal/database"
	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"github.com/alhazenlabs/procure-hub/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with the provided information.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User registration request"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {object} object "Bad Request"
// @Failure 409 {object} object "Conflict"
// @Router /users/v1/signup [post]
func RegisterUser(c *gin.Context) {
	// Initialize a new User instance
	var user models.User

	logger.Info("starting the request")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user input
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Save the user to the database using GORM
	db := database.GetDBWithCtx(c.Request.Context())
	db.Create(&user)

	// You can return a success response here if needed
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
