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

func TestGetDocumentsByUserID(t *testing.T) {
	database.Connect()

	user := models.User{Email: "getdocs@test.com", Password: "123"}
	database.DB.Create(&user)

	doc1 := models.Document{
		PublicID: "doc1-public",
		Title:    "Doc 1",
		UserID:   user.ID,
	}
	database.DB.Create(&doc1)

	doc2 := models.Document{
		PublicID: "doc2-public",
		Title:    "Doc 2",
		UserID:   user.ID,
	}
	database.DB.Create(&doc2)

	user2 := models.User{Email: "getdocs2@test.com", Password: "123"}
	database.DB.Create(&user2)
	doc3 := models.Document{
		PublicID: "doc3-public",
		Title:    "Doc 3",
		UserID:   user2.ID,
	}
	database.DB.Create(&doc3)

	documents, err := services.GetDocumentsByUserID(user.ID)

	assert.NoError(t, err)
	assert.Len(t, documents, 2, "Deveria retornar 2 documentos para o usuário")
	assert.Equal(t, "Doc 1", documents[0].Title)
	assert.Equal(t, "Doc 2", documents[1].Title)
}
