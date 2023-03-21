package shared

import (
	"context"

	"github.com/si-bas/go-rest-geospatial/domain/model"
)

func GetContextValueAsString(ctx context.Context, key string) string {
	val, ok := ctx.Value(key).(string)
	if ok {
		return val
	}

	return ""
}

func ChunkGeospatialData(data []model.Geospatial, chunkSize int) [][]model.Geospatial {
	var chunks [][]model.Geospatial
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}
