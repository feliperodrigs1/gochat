package services

import (
	"encoding/json"
	"errors"
	"gochat/internal/database"
	"gochat/internal/models"
	"io"

	"github.com/google/uuid"
)

func ProcessAndSaveDocument(userID uint, filename string, fileReader io.Reader) (string, int, error) {
	contentBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return "", 0, errors.New("failed to read file content")
	}

	text := string(contentBytes)
	chunks := SplitTextIntoChunks(text, 500)

	tx := database.DB.Begin()

	doc := models.Document{
		PublicID: uuid.NewString(),
		Title:    filename,
		UserID:   userID,
	}

	if err := tx.Create(&doc).Error; err != nil {
		tx.Rollback()
		return "", 0, errors.New("failed to create document")
	}

	for _, chunk := range chunks {
		embedding, err := GenerateEmbedding(chunk)
		if err != nil {
			tx.Rollback()
			return "", 0, err
		}

		embeddingJSON, _ := json.Marshal(embedding)

		chunkModel := models.Chunk{
			DocumentID: doc.ID,
			Content:    chunk,
			Embedding:  string(embeddingJSON),
		}

		if err := tx.Create(&chunkModel).Error; err != nil {
			tx.Rollback()
			return "", 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", 0, err
	}

	return doc.PublicID, len(chunks), nil
}

func SplitTextIntoChunks(text string, chunkSize int) []string {
	var chunks []string

	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize

		if end > len(runes) {
			end = len(runes)
		}

		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}

func GetDocumentsByUserID(userID uint) ([]models.Document, error) {
	var documents []models.Document
	if err := database.DB.Where("user_id = ?", userID).Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}
