package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func AskOpenAI(context, question string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	prompt := `
	You are a strict assistant.

	Answer ONLY using the provided context.
	If the answer is not in the context, say:
	"I don't know based on the provided document."

	DO NOT:
	- invent information
	- use prior knowledge
	- guess

	DO:
	- use the same language as the question
	- be concise
	- be direct
	- be clear

	Context:
	` + context + `

	Question:
	` + question

	body := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("invalid or empty choices returned from openai")
	}

	return result.Choices[0].Message.Content, nil
}
