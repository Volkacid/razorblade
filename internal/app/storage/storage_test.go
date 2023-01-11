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

func TestStorageGetValuesByID(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		UserID string
		err    error
	}{
		{
			name:   "Trying to get values with correct UserID",
			key:    "correct",
			UserID: "correctid",
			err:    nil,
		},
		{
			name:   "Trying to get values with wrong UserID",
			key:    "wrong",
			UserID: "wrongid",
			err:    NotFoundError(),
		},
	}
	db := CreateTestStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.SaveValue(context.Background(), tt.key, "https://ya.ru", "correctid")
			_, err := db.GetValuesByID(context.Background(), tt.UserID)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestStorageFindDuplicate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		err   error
	}{
		{
			name:  "Trying to get value that already in map",
			value: "https://ya.ru",
			err:   nil,
		},
		{
			name:  "Trying to get value that isn't stored",
			value: "wrong",
			err:   NotFoundError(),
		},
	}
	db := CreateTestStorage()
	db.SaveValue(context.Background(), "correct", tests[0].value, "someid")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := db.FindDuplicate(context.Background(), tt.value)
			assert.Equal(t, tt.err, err)
		})
	}
}
