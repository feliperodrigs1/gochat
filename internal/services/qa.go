package services

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"gochat/internal/database"
	"gochat/internal/models"
)

type ScoredChunk struct {
	Content string
	Score   float64
}

func AnswerWithConversation(userID uint, documentID, question, externalID string) (string, error) {
	normalizedQuestion := strings.ToLower(strings.TrimSpace(question))

	var doc models.Document
	if err := database.DB.
		Where("public_id = ? AND user_id = ?", documentID, userID).
		First(&doc).Error; err != nil {

		return "", errors.New("document not found")
	}

	var conv models.Conversation

	err := database.DB.
		Where("document_id = ? AND external_id = ?", doc.ID, externalID).
		First(&conv).Error

	if err != nil {
		conv = models.Conversation{
			DocumentID: doc.ID,
			ExternalID: externalID,
		}
		database.DB.Create(&conv)
	}

	database.DB.Create(&models.Message{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        normalizedQuestion,
	})

	var messages []models.Message

	database.DB.
		Where("conversation_id = ?", conv.ID).
		Order("id ASC").
		Limit(10).
		Find(&messages)

	history := ""
	for _, m := range messages {
		history += m.Role + ": " + m.Content + "\n"
	}

	questionEmbedding, err := GenerateEmbedding(normalizedQuestion)
	if err != nil {
		return "", err
	}

	var chunks []models.Chunk
	database.DB.
		Where("document_id = ?", doc.ID).
		Limit(50).
		Find(&chunks)

	var scored []ScoredChunk

	for _, ch := range chunks {
		var emb []float64
		json.Unmarshal([]byte(ch.Embedding), &emb)

		score := CosineSimilarity(questionEmbedding, emb)

		scored = append(scored, ScoredChunk{
			Content: ch.Content,
			Score:   score,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	topK := 5
	if len(scored) < topK {
		topK = len(scored)
	}

	context := ""
	for i := 0; i < topK; i++ {
		context += scored[i].Content + "\n"
	}

	fullPrompt := "Answer ONLY using the context.\n\n" +
		"Context:\n" + context + "\n\n" +
		"Conversation:\n" + history + "\n\n" +
		"Question:\n" + normalizedQuestion

	answer, err := AskOpenAI(fullPrompt, "")
	if err != nil {
		return "", err
	}

	database.DB.Create(&models.Message{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        answer,
	})

	return answer, nil
}
