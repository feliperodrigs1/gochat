package services_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"gochat/internal/services"
)

func TestGenerateEmbedding_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"data": []map[string]interface{}{
				{"embedding": []float64{0.1, 0.2, 0.3}},
			},
		}),
	)

	emb, err := services.GenerateEmbedding("test text")

	assert.NoError(t, err)
	assert.NotNil(t, emb)
	assert.Equal(t, []float64{0.1, 0.2, 0.3}, emb)
}

func TestGenerateEmbedding_HTTPError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
		httpmock.NewStringResponder(500, "Internal Server Error"),
	)

	emb, err := services.GenerateEmbedding("test error text")

	assert.Error(t, err)
	assert.Equal(t, "no embedding returned", err.Error())
	assert.Nil(t, emb)
}

func TestGenerateEmbedding_NoData(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"data": []interface{}{},
		}),
	)

	emb, err := services.GenerateEmbedding("test empty text")

	assert.Error(t, err)
	assert.Equal(t, "no embedding returned", err.Error())
	assert.Nil(t, emb)
}
