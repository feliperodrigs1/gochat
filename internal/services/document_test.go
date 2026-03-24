package services_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"gochat/internal/database"
	"gochat/internal/models"
	"gochat/internal/services"
)

func TestSplitTextIntoChunks(t *testing.T) {
	text := "1234567890"

	chunks := services.SplitTextIntoChunks(text, 3)

	assert.Len(t, chunks, 4, "deveria ter quebrado em 4 partes")
	assert.Equal(t, "123", chunks[0])
	assert.Equal(t, "456", chunks[1])
	assert.Equal(t, "789", chunks[2])
	assert.Equal(t, "0", chunks[3])
}

func TestProcessAndSaveDocument(t *testing.T) {
	os.Setenv("ENV", "test")
	database.Connect()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"data": []map[string]interface{}{
				{"embedding": []float64{0.1, 0.2}},
			},
		}),
	)

	user := models.User{Email: "svc@test.com", Password: "123"}
	database.DB.Create(&user)

	fileContent := bytes.NewBufferString("Este é o meu documento de teste para testar a gravação.")

	docID, totalChunks, err := services.ProcessAndSaveDocument(user.ID, "test.txt", fileContent)

	assert.NoError(t, err)
	assert.NotEmpty(t, docID, "O PublicID do documento não devia estar vazio")
	assert.Equal(t, 1, totalChunks, "O texto é pequeno, deveria gerar apenas 1 chunk")

	var doc models.Document
	database.DB.First(&doc, "public_id = ?", docID)
	assert.Equal(t, "test.txt", doc.Title)
	assert.Equal(t, user.ID, doc.UserID)
}
