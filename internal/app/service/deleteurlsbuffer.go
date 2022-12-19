package service

import (
	"context"
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"strings"
	"time"
)

type URLsDeleteBuffer struct {
	deleteChan chan string
	keysBuffer []string
	db         storage.Storage
}

func NewDeleteBuffer(db storage.Storage) *URLsDeleteBuffer {
	buf := &URLsDeleteBuffer{deleteChan: make(chan string), keysBuffer: make([]string, 0, config.DeleteBufferSize), db: db}
	buf.newDeleteWorker()
	return buf
}

func (buffer *URLsDeleteBuffer) newDeleteWorker() {
	go func() {
		bufTicker := time.NewTicker(config.DeleteBufferClearTime * time.Second)
		for {
			select {
			case <-bufTicker.C:
				buffer.cleanBuffer()
			case incomeKey := <-buffer.deleteChan:
				buffer.keysBuffer = append(buffer.keysBuffer, incomeKey)
				if len(buffer.keysBuffer) == config.DeleteBufferSize {
					buffer.cleanBuffer()
				}
			}
		}
	}()
}

func (buffer *URLsDeleteBuffer) AddKeys(keys []string, userID string) {
	userURLs, err := buffer.db.GetValuesByID(context.Background(), userID)
	if err != nil {
		fmt.Println("Cannot get user URLs: ", err)
		return
	}
	var userKeys string
	for _, userKey := range userURLs {
		userKeys += userKey.ShortURL
	}
	for _, key := range keys {
		if strings.Contains(userKeys, key) {
			buffer.deleteChan <- key
		}
	}
}

func (buffer *URLsDeleteBuffer) cleanBuffer() {
	if len(buffer.keysBuffer) != 0 {
		buffer.db.DeleteURLs(buffer.keysBuffer)
		buffer.keysBuffer = buffer.keysBuffer[:0]
	}
}
