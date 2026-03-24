package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gochat/internal/services"
)

func Ask(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body struct {
		DocumentID string `json:"document_id"`
		Question   string `json:"question"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if body.DocumentID == "" || body.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID and question are required"})
		return
	}

	answer, err := services.AnswerQuestion(userID, body.DocumentID, body.Question)
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
