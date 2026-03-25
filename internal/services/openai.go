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
	You are an assistant that answers questions based ONLY on the provided context.

	STRICT RULES:
	- Use ONLY the information from the context
	- If the answer is not in the context, respond in the same language as the question saying:
	"I don't know based on the provided document."
	- Do NOT invent or assume anything
	- Do NOT use prior knowledge

	STYLE:
	- Use the same language as the question
	- Be concise
	- Be clear
	- Prefer exact phrases from the context

	CONTEXT:
	` + context + `

	CONVERSATION + QUESTION:
	` + question + `
	`

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
	req.Header.Set("Authorization", "Bearer " + apiKey)

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

func AskOpenAIRewrite(history, question string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	prompt := `
	Rewrite the user's question to be fully self-contained.

	Use the conversation history for context.

	DO NOT answer the question.
	ONLY rewrite it.

	Conversation:
	` + history + `

	Question:
	` + question + `

	Rewritten question:
	`

	body := map[string]interface{}{
		"model":       "gpt-4o-mini",
		"temperature": 0,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

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
			}
		}
	}

	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Choices) == 0 {
		return question, nil
	}

	return result.Choices[0].Message.Content, nil
}
