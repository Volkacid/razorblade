package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestDeleteUserURLs(t *testing.T) {
	t.Run("Trying to delete URLs", func(t *testing.T) {
		db := storage.CreateTestStorage()
		keys := []string{"somekey", "otherkey"}
		marshalledKeys, _ := json.Marshal(keys)
		userID := "someuserid"
		for _, key := range keys {
			db.SaveValue(context.Background(), key, "https://ya.ru", userID)
		}
		reader := bytes.NewReader(marshalledKeys)
		request := TestRequest("/api/user/urls", http.MethodDelete, reader, db, userID)
		response := request.Result()
		defer response.Body.Close()
		assert.Equal(t, http.StatusAccepted, response.StatusCode)
		time.Sleep(1 * time.Second)
		for _, key := range keys {
			request = TestRequest("/"+key, http.MethodGet, nil, db, userID)
			response = request.Result()
			response.Body.Close()
			assert.Equal(t, http.StatusGone, response.StatusCode)
		}
	})

}
