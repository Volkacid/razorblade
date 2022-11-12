package storage

import (
	"errors"
	"sync"
)

type Storage struct {
	db map[string]string
}

var mutex = &sync.RWMutex{}

func CreateStorage() *Storage {
	return &Storage{db: make(map[string]string)}
}

func (storage *Storage) GetValue(key string) (string, error) {
	mutex.RLock()
	value, ok := storage.db[key]
	mutex.RUnlock()
	if !ok {
		return "", errors.New("value does not exist")
	} else {
		return value, nil
	}
}

func (storage *Storage) SaveValue(key string, value string) {
	mutex.Lock()
	storage.db[key] = value
	mutex.Unlock()
}
