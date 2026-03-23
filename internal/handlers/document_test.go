package handlers_test

import (
	"bytes"
	"mime/multipart"
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

func setupDocumentRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENV", "test")
	database.Connect()

	database.DB.Create(&models.User{Email: "handlerdoc@test.com", Password: "123"})

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.POST("/documents", handlers.CreateDocument)

	return r
}

func createMultipartRequest(t *testing.T, filename string, content string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)

	part.Write([]byte(content))
	writer.Close()

	req, err := http.NewRequest("POST", "/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func TestCreateDocumentSuccess(t *testing.T) {
	r := setupDocumentRouter()

	req, err := createMultipartRequest(t, "documento.txt", "Exemplo de conteúdo salvo via arquivo")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.Contains(t, w.Body.String(), "Document uploaded successfully")
	assert.Contains(t, w.Body.String(), "documentID")
}

func TestCreateDocumentWithoutFile(t *testing.T) {
	r := setupDocumentRouter()

	req, _ := http.NewRequest("POST", "/documents", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
