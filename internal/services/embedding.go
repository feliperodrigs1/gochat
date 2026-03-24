package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func GenerateEmbedding(text string) ([]float64, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	body := map[string]interface{}{
		"input": text,
		"model": "text-embedding-3-small",
	}

	jsonBody, err := json.Marshal(body)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonBody))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbeddingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if  len(result.Data) == 0 {
		return nil, errors.New("no embedding returned")
	}

	return result.Data[0].Embedding, nil
}
