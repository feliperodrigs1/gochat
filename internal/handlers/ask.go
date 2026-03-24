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
		ExternalID string   `json:"external_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if body.DocumentID == "" || body.Question == "" || body.ExternalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	answer, err := services.AnswerWithConversation(userID, body.DocumentID, body.Question, body.ExternalID)
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
