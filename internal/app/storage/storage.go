package storage

import (
	"bufio"
	"errors"
	"github.com/Volkacid/razorblade/internal/app/config"
	"os"
	"strings"
	"sync"
)

type Storage struct {
	db map[string]string
}

var mutex = &sync.RWMutex{}
var storageFilePath string
var storageFileExist bool

func CreateStorage(byFile bool) *Storage {
	if byFile {
		storageFilePath = config.GetServerConfig().StorageFile
		storageFileExist = storageFilePath != ""
		return &Storage{}
	} else {
		return &Storage{db: make(map[string]string)}
	}
}

func (storage *Storage) GetValue(key string) (string, error) {
	if storageFileExist {
		db, err := os.OpenFile(storageFilePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		foundValue := ""
		mutex.RLock()
		scanner := bufio.NewScanner(db)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), key) {
				_, foundValue, _ = strings.Cut(scanner.Text(), ":-:")
				break
			}
		}
		mutex.RUnlock()
		if err = db.Close(); err != nil {
			panic(err)
		}
		if foundValue != "" {
			return foundValue, nil
		} else {
			return "", errors.New("value not found")
		}
	} else {
		mutex.RLock()
		value, ok := storage.db[key]
		mutex.RUnlock()
		if !ok {
			return "", errors.New("value not found")
		} else {
			return value, nil
		}
	}
}

func (storage *Storage) SaveValue(key string, value string) {
	if storageFileExist {
		db, err := os.OpenFile(storageFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		_, err = db.WriteString(key + ":-:" + value + "\n")
		if err != nil {
			panic(err)
		}
		mutex.Unlock()
		if err = db.Close(); err != nil {
			panic(err)
		}
	} else {
		mutex.Lock()
		storage.db[key] = value
		mutex.Unlock()
	}
}
