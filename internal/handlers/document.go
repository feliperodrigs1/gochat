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

func GetDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	documents, err := services.GetDocumentsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	type DocumentResponse struct {
		PublicID  string `json:"public_id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
	}

	var response []DocumentResponse
	for _, doc := range documents {
		response = append(response, DocumentResponse{
			PublicID:  doc.PublicID,
			Title:     doc.Title,
			CreatedAt: doc.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, response)
}
