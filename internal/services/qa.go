package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"gochat/internal/cache"
	"gochat/internal/database"
	"gochat/internal/models"
)

type ScoredChunk struct {
	Content string
	Score   float64
}

var ctx = context.Background()

func AnswerWithConversation(userID uint, documentID, question, externalID string) (string, error) {
	normalizedQuestion := strings.ToLower(strings.TrimSpace(question))

	cacheKey := "qa:" + documentID + ":" + normalizedQuestion
	cached, err := cache.Client.Get(ctx, cacheKey).Result()

	if err == nil && cached != "" {
		return cached, nil
	}

	doc, err := fetchDocument(userID, documentID)
	if err != nil {
		return "", err
	}

	conv, err := getOrCreateConversation(doc.ID, externalID)
	if err != nil {
		return "", err
	}

	if err := saveUserMessage(conv.ID, normalizedQuestion); err != nil {
		return "", err
	}

	history, err := buildConversationHistory(conv.ID)
	if err != nil {
		return "", err
	}

	answer, err := generateAnswer(doc.ID, normalizedQuestion, history)

	if err == nil && !isBadAnswer(answer) {
		saveAssistantMessage(conv.ID, answer)

		cache.Client.Set(ctx, cacheKey, answer, 10*time.Minute)

		return answer, nil
	}

	rewritten, err := RewriteQuestion(history, normalizedQuestion)
	if err != nil {
		log.Println("error at rewrite, using original")
		rewritten = normalizedQuestion
	}

	answer, err = generateAnswer(doc.ID, rewritten, history)
	if err != nil {
		return "", err
	}

	saveAssistantMessage(conv.ID, answer)

	cache.Client.Set(ctx, cacheKey, answer, 10*time.Minute)

	return answer, nil
}

func fetchDocument(userID uint, documentID string) (models.Document, error) {
	var doc models.Document
	if err := database.DB.
		Where("public_id = ? AND user_id = ?", documentID, userID).
		First(&doc).Error; err != nil {
		return doc, errors.New("document not found")
	}

	return doc, nil
}

func getOrCreateConversation(docID uint, externalID string) (models.Conversation, error) {
	var conv models.Conversation
	err := database.DB.
		Where("document_id = ? AND external_id = ?", docID, externalID).
		First(&conv).Error

	if err != nil {
		conv = models.Conversation{
			DocumentID: docID,
			ExternalID: externalID,
		}
		if err := database.DB.Create(&conv).Error; err != nil {
			return conv, err
		}
	}
	return conv, nil
}

func saveUserMessage(conversationID uint, content string) error {
	return database.DB.Create(&models.Message{
		ConversationID: conversationID,
		Role:           "user",
		Content:        content,
	}).Error
}

func buildConversationHistory(conversationID uint) (string, error) {
	var messagesDesc []models.Message

	err := database.DB.
		Where("conversation_id = ?", conversationID).
		Order("id DESC").
		Limit(6).
		Find(&messagesDesc).Error
	if err != nil {
		return "", err
	}

	var messages []models.Message
	for i := len(messagesDesc) - 1; i >= 0; i-- {
		messages = append(messages, messagesDesc[i])
	}

	history := ""
	for _, m := range messages {
		history += m.Role + ": " + m.Content + "\n"
	}

	return history, nil
}

func generateAnswer(documentID uint, question, history string) (string, error) {
	questionEmbedding, err := GenerateEmbedding(question)
	if err != nil {
		return "", err
	}

	chunks, err := fetchDocumentChunks(documentID)
	if err != nil {
		return "", err
	}

	topChunks := getTopScoredChunks(chunks, question, questionEmbedding, 5)

	if len(topChunks) == 0 {
		return "I don't know based on the provided document.", nil
	}

	context := buildContextFromChunks(topChunks)

	return AskOpenAI(context, history+"\nQuestion: "+question)
}

func fetchDocumentChunks(documentID uint) ([]models.Chunk, error) {
	var chunks []models.Chunk
	err := database.DB.
		Where("document_id = ?", documentID).
		Limit(50).
		Find(&chunks).Error
	return chunks, err
}

func getTopScoredChunks(chunks []models.Chunk, question string, questionEmbedding []float64, topK int) []ScoredChunk {
	var scored []ScoredChunk
	words := strings.Fields(strings.ToLower(question))

	for _, ch := range chunks {
		var emb []float64
		if err := json.Unmarshal([]byte(ch.Embedding), &emb); err != nil {
			continue
		}

		score := CosineSimilarity(questionEmbedding, emb)
		contentLower := strings.ToLower(ch.Content)

		for _, word := range words {
			if len(word) < 4 {
				continue
			}
			if strings.Contains(contentLower, word) {
				score += 0.15
			}
		}

		if strings.Contains(contentLower, question) {
			score += 0.2
		}

		scored = append(scored, ScoredChunk{
			Content: ch.Content,
			Score:   score,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if len(scored) == 0 {
		return []ScoredChunk{}
	}

	if len(scored) < topK {
		topK = len(scored)
	}

	return scored[:topK]
}

func buildContextFromChunks(chunks []ScoredChunk) string {
	context := ""
	for i, t := range chunks {
		context += "Source " + strconv.Itoa(i+1) + ":\n" + t.Content + "\n\n"
	}
	return context
}

func isBadAnswer(answer string) bool {
	a := strings.ToLower(answer)

	return strings.Contains(a, "i don't know") ||
		strings.Contains(a, "não sei") ||
		strings.Contains(a, "não encontrado") ||
		len(a) < 10
}

func saveAssistantMessage(conversationID uint, answer string) {
	database.DB.Create(&models.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        answer,
	})
}

func RewriteQuestion(history, question string) (string, error) {
	return AskOpenAIRewrite(history, question)
}
