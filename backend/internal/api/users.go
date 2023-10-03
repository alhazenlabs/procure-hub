package api

import (
	"context"
	"net/http"

	"github.com/alhazenlabs/procure-hub/backend/internal/database"
	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"github.com/alhazenlabs/procure-hub/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Debug("recieved malformed request", zap.Any("error", err))
		return
	}
	logger.Info("recieved request to create", zap.String("user", user.Email))

	// Validate user input
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error("failed to vaidate request", zap.String("user", user.Email), zap.Any("error", err))
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errStr := "Failed to hash the password"
		c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
		logger.Error(errStr, zap.String("user", user.Email), zap.Any("error", err))
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Save the user to the database using GORM
	tx, commitFunc := database.BeginTransaction(c.Request.Context())
	defer commitFunc()

	// Logic below to check if the length of users is greater than 1
	primaryOwner := models.User{}
	query := tx.Find(&primaryOwner)
	if query.Error != nil {
		logger.Error(query.Error.Error())
	}
	if query.RowsAffected > 0 {
		// Check if the PrimaryOwnerId Exists in the table and is primary owner
		tx.Where("id = ?", user.PrimaryOwnerID).First(&primaryOwner)
		if primaryOwner.ID != user.PrimaryOwnerID {
			errStr := "Primary owner doesn't exists"
			c.JSON(http.StatusBadRequest, gin.H{"error": errStr})
			logger.Error(errStr, zap.String("user", user.Email))
			return
		}
		user.UserType = models.Client
	}

	if err := tx.Create(&user).Error; err != nil {
		errStr := "Failed creating user"
		c.JSON(http.StatusConflict, gin.H{"error": errStr})
		logger.Error(errStr, zap.String("user", user.Email), zap.Any("error", err))
		return
	}
	tx.Commit()
	logger.Info("created", zap.String("user", user.Email))

	// You can return a success response here if needed
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// @Summary User Login
// @Description Log in with the provided email and password.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body api.LoginHandler.LoginData true "User login request"
// @Success 200 {string} string "Login successful"
// @Failure 400 {object} object "Bad Request"
// @Failure 401 {object} object "Unauthorized"
// @Router /users/v1/login [post]
func LoginHandler(c *gin.Context) {
	type LoginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	loginData := LoginData{}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Debug("recieved malformed request")
		return
	}
	logger.Info("recieved a login request", zap.String("email", loginData.Email))

	user, err := findUserByEmail(c.Request.Context(), loginData.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		logger.Error("Invalid email", zap.String("user", user.Email))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		logger.Error("Invalid Password", zap.String("user", user.Email))
		return
	}

	logger.Info("login successful", zap.String("email", loginData.Email))
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func findUserByEmail(context context.Context, email string) (*models.User, error) {
	// Replace this with your logic to retrieve a user by email from the database.
	// This is a placeholder function, and you should implement your own logic.
	// Save the user to the database using GORM
	tx, commitFunc := database.BeginTransaction(context)
	defer commitFunc()

	user := models.User{}
	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	tx.Commit()

	return &user, nil
}
