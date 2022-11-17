package storage

import (
	"bufio"
	"errors"
	"fmt"
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
	fmt.Println("IS EXIST", storageFileExist)
	if storageFileExist {
		db, err := os.OpenFile(storageFilePath, os.O_RDONLY, 0444)
		defer func(db *os.File) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)
		if err != nil {
			return "", err
		}
		mutex.RLock()
		scanner := bufio.NewScanner(db)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), key) {
				_, foundValue, _ := strings.Cut(scanner.Text(), ":-:")
				return foundValue, nil
			}
		}
		mutex.RUnlock()
		return "", errors.New("value not found")
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
		defer func(db *os.File) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		db.WriteString(key + ":-:" + value + "\n")
		mutex.Unlock()
	} else {
		mutex.Lock()
		storage.db[key] = value
		mutex.Unlock()
	}
}
