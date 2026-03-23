package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gochat/internal/database"
	"gochat/internal/handlers"
	"gochat/internal/models"
	"gochat/internal/services"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	os.Setenv("ENV", "test")

	database.Connect()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	return r
}

func TestRegister(t *testing.T) {
	r := setupRouter()

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLoginSuccess(t *testing.T) {
	r := setupRouter()

	hash, _ := services.HashPassword("123456")

	user := models.User{
		Email:    "login@test.com",
		Password: hash,
	}

	database.DB.Create(&user)

	body := map[string]string{
		"email":    "login@test.com",
		"password": "123456",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLoginInvalid(t *testing.T) {
	r := setupRouter()

	body := map[string]string{
		"email":    "naoexiste@test.com",
		"password": "123456",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}
