package test

import (
	"context"
	"reflect"
	"testing"

	"github.com/si-bas/go-rest-geospatial/domain/model"
	"github.com/si-bas/go-rest-geospatial/shared"
)

func TestGetContextValueAsString(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")

	result := shared.GetContextValueAsString(ctx, "key")

	if result != "value" {
		t.Errorf("GetContextValueAsString returned unexpected result, got: %s, want: %s.", result, "value")
	}
}

func TestChunkGeospatialData(t *testing.T) {
	data := []model.Geospatial{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}

	expectedChunks := [][]model.Geospatial{
		{
			{ID: 1},
			{ID: 2},
		},
		{
			{ID: 3},
		},
	}

	chunkSize := 2

	chunks := shared.ChunkGeospatialData(data, chunkSize)

	if !reflect.DeepEqual(chunks, expectedChunks) {
		t.Errorf("expected chunks: %v, got: %v", expectedChunks, chunks)
	}
}
