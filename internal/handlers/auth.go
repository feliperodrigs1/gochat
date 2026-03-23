package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gochat/internal/models"
	"gochat/internal/database"
	"gochat/internal/services"
)

type AuthInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input AuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hash, _ := services.HashPassword(input.Password)

	user := models.User{
		Email:    input.Email,
		Password: hash,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func Login(c *gin.Context) {
	var input AuthInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := services.GenerateToken(user.ID)

	c.JSON(http.StatusOK, gin.H{"token": token})
}