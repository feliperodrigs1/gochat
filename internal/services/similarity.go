package services

import "math"

func CosineSimilarity(a, b []float64) float64 {
	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
