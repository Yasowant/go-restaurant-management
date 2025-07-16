package controllers

import (
	"log"
	"net/http"
	"restaurant-app/models"
	"restaurant-app/services"
	"restaurant-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ✅ REGISTER CONTROLLER
func Register(c *gin.Context) {
	var user models.User

	// Bind JSON input
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validate input struct
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save user to database
	err := services.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send welcome email and log result
	if err := utils.SendEmail(
		user.Email,
		"Welcome to Restaurant App!",
		"<h1>Hi "+user.FullName+",</h1><p>Thanks for registering 🍽️</p>",
	); err != nil {
		log.Println("❌ Failed to send email:", err)
	} else {
		log.Println("✅ Email sent successfully to", user.Email)
	}

	// Send response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// ✅ LOGIN CONTROLLER
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// Bind input
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validate struct
	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate and get JWT token
	token, err := services.LoginUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Respond with token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := services.UpdateUserProfile(userID, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var body struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ChangePassword(userID, body.OldPassword, body.NewPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
