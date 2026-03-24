package handlers

import (
	"gochat/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	docPublicId, totalChunks, err := services.ProcessAndSaveDocument(userID.(uint), file.Filename, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Document uploaded successfully",
		"documentID": docPublicId,
		"chunks":     totalChunks,
	})
}
