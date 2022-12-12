package storage

import (
	"context"
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/config"
	"sync"
)

type Storage interface {
	GetValue(key string) (string, error)
	GetValuesByID(userID string) ([]UserURL, error)
	SaveValue(key string, value string, userID string) error
	BatchSave(ctx context.Context, values map[string]string, userID string) error
	FindDuplicate(value string) (string, error)
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var mutex = &sync.RWMutex{}

func CreateStorage() Storage {
	servConf := config.GetServerConfig()
	if CheckDBConnection() {
		err := InitializeDB()
		if err != nil {
			fmt.Printf("DB initialization error(%v). Storage created in map.", err)
			return NewMapDB()
		}
		fmt.Println("Storage created in DB")
		return NewDB()
	}

	if servConf.StorageFile != "" {
		fmt.Println("Storage created in file")
		return NewFileDB(servConf.StorageFile)
	}

	fmt.Println("Storage created in map")
	return NewMapDB()
}

func CreateTestStorage() Storage {
	return NewMapDB()
}
