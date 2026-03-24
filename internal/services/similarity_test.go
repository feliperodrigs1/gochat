package services_test

import (
	"math"
	"testing"

	"gochat/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	a := []float64{1, 0, 0}
	b := []float64{1, 0, 0}
	score := services.CosineSimilarity(a, b)
	assert.Equal(t, float64(1), score, "Vectors are identical, similarity should be 1")

	a2 := []float64{1, 0}
	b2 := []float64{0, 1}
	score2 := services.CosineSimilarity(a2, b2)
	assert.Equal(t, float64(0), score2, "Orthogonal vectors, similarity should be 0")

	a3 := []float64{1, 1}
	b3 := []float64{-1, -1}
	score3 := services.CosineSimilarity(a3, b3)

	assert.True(t, math.Abs(score3 - -1) < 0.0001, "Opposite vectors, similarity should be approx -1")
}
