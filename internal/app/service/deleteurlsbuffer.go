package service

import (
	"context"
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"time"
)

type URLsDeleteBuffer struct {
	deleteChan chan string
	keysBuffer []string
	db         storage.Storage
	ctx        context.Context
}

func NewDeleteBuffer(db storage.Storage, ctx context.Context) *URLsDeleteBuffer {
	buf := &URLsDeleteBuffer{deleteChan: make(chan string),
		keysBuffer: make([]string, 0, config.DeleteBufferSize),
		db:         db,
		ctx:        ctx}
	buf.newDeletionWorker()
	return buf
}

func (buffer *URLsDeleteBuffer) newDeletionWorker() {
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
			case <-buffer.ctx.Done():
				bufTicker.Stop()
				buffer.cleanBuffer()
			}
		}
	}()
}

func (buffer *URLsDeleteBuffer) AddKeys(keys []string, userID string) {
	userURLs, err := buffer.db.GetValuesByID(buffer.ctx, userID)
	if err != nil {
		fmt.Println("Cannot get user URLs: ", err)
		return
	}
	userKeys := make(map[string]bool, len(userURLs))
	for _, userVal := range userURLs {
		userKeys[userVal.ShortURL] = true
	}
	for _, key := range keys {
		if userKeys[key] {
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
