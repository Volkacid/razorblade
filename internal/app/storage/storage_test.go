package storage

import (
	"errors"
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
			err:  errors.New("value does not exist"),
		},
	}
	db := CreateStorage()
	db.SaveValue("testkey", "somevalue")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, resultErr := db.GetValue(tt.key)
			assert.Equal(t, tt.err, resultErr)
		})
	}
}
