package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
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
	r.GET("/documents", handlers.GetDocuments)

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

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"data": []map[string]interface{}{
				{"embedding": []float64{0.1, 0.2}},
			},
		}),
	)

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

func TestGetDocuments(t *testing.T) {
	r := setupDocumentRouter()

	var user models.User
	database.DB.First(&user, 1)

	doc1 := models.Document{PublicID: "doc1-pub", Title: "Doc 1", UserID: user.ID}
	database.DB.Create(&doc1)

	doc2 := models.Document{PublicID: "doc2-pub", Title: "Doc 2", UserID: user.ID}
	database.DB.Create(&doc2)

	req, _ := http.NewRequest("GET", "/documents", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, "doc1-pub", response[0]["public_id"])
	assert.Equal(t, "Doc 1", response[0]["title"])
	assert.Equal(t, "doc2-pub", response[1]["public_id"])
	assert.Equal(t, "Doc 2", response[1]["title"])
}
