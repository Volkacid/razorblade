package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageGetValue(t *testing.T) {
	tests := []struct {
		name string
		key  string
		err  error
	}{
		{
			name: "Trying to save and get original value",
			key:  "testkey",
			err:  nil,
		},
		{
			name: "Trying to get value that isn't stored",
			key:  "anothertestkey",
			err:  NotFoundError(),
		},
	}
	db := CreateTestStorage()
	ctx := context.Background()
	db.SaveValue(ctx, "testkey", "somevalue", "someid")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, resultErr := db.GetValue(ctx, tt.key)
			assert.Equal(t, tt.err, resultErr)
		})
	}
}
