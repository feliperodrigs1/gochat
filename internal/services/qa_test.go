package services_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"gochat/internal/database"
	"gochat/internal/models"
	"gochat/internal/services"
)

func TestAnswerQuestion_Success(t *testing.T) {
	os.Setenv("ENV", "test")
	database.Connect()

	database.DB.Exec("DELETE FROM chunks")
	database.DB.Exec("DELETE FROM documents")

	doc := models.Document{
		Title:    "Test Doc",
		PublicId: "qa-test-id",
		UserID:   99,
	}
	database.DB.Create(&doc)

	chunk := models.Chunk{
		DocumentID: doc.ID,
		Content:    "O gato mia",
		Embedding:  "[0.1, 0.2, 0.3]",
	}
	database.DB.Create(&chunk)

	chunk2 := models.Chunk{
		DocumentID: doc.ID,
		Content:    "O cachorro late",
		Embedding:  "[0.9, 0.8, 0.7]",
	}
	database.DB.Create(&chunk2)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"data": []map[string]interface{}{
				{"embedding": []float64{0.1, 0.2, 0.3}},
			},
		}),
	)

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "O gato mia sim",
					},
				},
			},
		}),
	)

	answer, err := services.AnswerQuestion(99, "qa-test-id", "o que o gato faz?")

	assert.NoError(t, err)
	assert.Equal(t, "O gato mia sim", answer)
}

func TestAnswerQuestion_DocumentNotFound(t *testing.T) {
	os.Setenv("ENV", "test")
	database.Connect()

	answer, err := services.AnswerQuestion(99, "unknown-id", "ola?")

	assert.Error(t, err)
	assert.Equal(t, "document not found", err.Error())
	assert.Equal(t, "", answer)
}
