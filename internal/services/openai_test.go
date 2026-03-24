package services_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"gochat/internal/services"
)

func TestAskOpenAI_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResponse := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": "A resposta da api",
				},
			},
		},
	}

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
		httpmock.NewJsonResponderOrPanic(200, mockResponse),
	)

	answer, err := services.AskOpenAI("texto de contexto", "qual a pergunta?")

	assert.NoError(t, err)
	assert.Equal(t, "A resposta da api", answer)
}

func TestAskOpenAI_InvalidResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
		httpmock.NewStringResponder(500, "Internal Server Error"),
	)

	answer, err := services.AskOpenAI("contexto", "pergunta")

	assert.Error(t, err)
	assert.Equal(t, "", answer)
}
