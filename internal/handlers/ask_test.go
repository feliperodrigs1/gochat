package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gochat/internal/database"
	"gochat/internal/handlers"
	"gochat/internal/models"
)

func setupAskRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENV", "test")
	database.Connect()

	database.DB.Create(&models.User{Email: "handlerask@test.com", Password: "123"})

	database.DB.Create(&models.Document{
		Title:    "Test Doc",
		PublicId: "test-doc-id",
		UserID:   1,
	})

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.POST("/ask", handlers.Ask)

	return r
}

func TestAskSuccess(t *testing.T) {
	r := setupAskRouter()

	body := map[string]string{
		"document_id": "test-doc-id",
		"question":    "Qual o significado da vida?",
		"external_id": "client_123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/ask", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestAskInvalidJSON(t *testing.T) {
	r := setupAskRouter()

	req, _ := http.NewRequest("POST", "/ask", bytes.NewBuffer([]byte(`{"invalid": "json`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAskMissingFields(t *testing.T) {
	r := setupAskRouter()

	body := map[string]string{
		"document_id": "",
		"question":    "",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/ask", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAskDocumentNotFound(t *testing.T) {
	r := setupAskRouter()

	body := map[string]string{
		"document_id": "non-existent-id",
		"question":    "Qual o significado da vida?",
		"external_id": "client_123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/ask", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
