package services

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"gochat/internal/database"
	"gochat/internal/models"
)

type ScoredChunk struct {
	Content string
	Score   float64
}

func AnswerQuestion(userID uint, documentID, question string) (string, error) {
	normalizedQuestion := strings.ToLower(strings.TrimSpace(question))

	var doc models.Document
	if err := database.DB.Where("public_id = ? AND user_id = ?", documentID, userID).First(&doc).Error; err != nil {
		return "", errors.New("document not found")
	}

	questionEmbedding, err := GenerateEmbedding(normalizedQuestion)
	if err != nil {
		return "", errors.New("failed to generate embedding")
	}

	var questions []models.Question

	database.DB.
		Where("document_id = ?", doc.ID).
		Limit(20).
		Find(&questions)

	bestScore := 0.0
	bestAnswer := ""

	for _, q := range questions {
		var emb []float64
		if err := json.Unmarshal([]byte(q.Embedding), &emb); err != nil {
			continue
		}

		score := CosineSimilarity(questionEmbedding, emb)

		if score > bestScore {
			bestScore = score
			bestAnswer = q.Answer
		}
	}

	if bestScore > 0.85 {
		return bestAnswer, nil
	}

	var chunks []models.Chunk
	if err := database.DB.Where("document_id = ?", doc.ID).Limit(50).Find(&chunks).Error; err != nil {
		return "", errors.New("failed to retrieve document chunks")
	}

	var scored []ScoredChunk

	words := strings.Fields(strings.ToLower(normalizedQuestion))

	for _, ch := range chunks {
		var emb []float64
		if err := json.Unmarshal([]byte(ch.Embedding), &emb); err != nil {
			continue
		}

		score := CosineSimilarity(questionEmbedding, emb)
		contentLower := strings.ToLower(ch.Content)

		for _, word := range words {
			if strings.Contains(contentLower, word) {
				score += 0.02
			}
		}

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

	if topK == 0 {
		return "No relevant information found in the document.", nil
	}

	top := scored[:topK]
	context := ""
	for i, t := range top {
		context += "Chunk " + strconv.Itoa(i+1) + ":\n" + t.Content + "\n\n"
	}

	answer, err := AskOpenAI(context, question)
	if err != nil {
		return "", errors.New("failed to get answer from open ai")
	}

	embeddingJSON, _ := json.Marshal(questionEmbedding)

	database.DB.Create(&models.Question{
		DocumentID: doc.ID,
		Question:   normalizedQuestion,
		Answer:     answer,
		Embedding:  string(embeddingJSON),
	})

	return answer, nil
}
