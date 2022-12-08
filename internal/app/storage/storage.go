package storage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/jackc/pgx/v5"
	"os"
	"strings"
	"sync"
	"time"
)

type Storage struct {
	connType        int
	dbConn          pgx.Conn
	storageFilePath string
	dbMap           map[string]string
}

const (
	ByDB = iota
	ByFile
	byMap
)

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var mutex = &sync.RWMutex{}

func CreateStorage() *Storage {
	servConf := config.GetServerConfig()
	if CheckDBConnection() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		dbConn, err := pgx.Connect(ctx, servConf.DBAddress)
		if err != nil {
			panic(err)
		}
		err = InitializeDB(dbConn)
		if err != nil {
			panic(err)
		}
		fmt.Println("Storage created in DB")
		return &Storage{connType: ByDB, dbConn: *dbConn}
	}
	if servConf.StorageFile != "" {
		fmt.Println("Storage created in file")
		return &Storage{connType: ByFile, storageFilePath: servConf.StorageFile}
	}
	fmt.Println("Storage created in map")
	return &Storage{connType: byMap, dbMap: make(map[string]string)}
}

func CreateTestStorage() *Storage {
	return &Storage{connType: byMap, dbMap: make(map[string]string)}
}

func (storage *Storage) GetValue(key string) (string, error) {
	switch storage.connType {
	case ByDB:
		var value string
		err := storage.dbConn.QueryRow(context.Background(), `SELECT original FROM urls WHERE short=$1`, key).Scan(&value)
		return value, err
	case ByFile:
		db, err := os.OpenFile(storage.storageFilePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		foundValue := ""
		mutex.RLock()
		scanner := bufio.NewScanner(db)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), key) {
				_, foundValue, _ = strings.Cut(scanner.Text(), ":-:")
				foundValue, _, _ = strings.Cut(foundValue, ":_:")
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
	case byMap:
		mutex.RLock()
		value, ok := storage.dbMap[key]
		mutex.RUnlock()
		if !ok {
			return "", errors.New("value not found")
		} else {
			return value, nil
		}
	}
	return "", nil
}

func (storage *Storage) GetValuesByID(userID string) ([]UserURLs, error) {
	var foundValues []UserURLs

	switch storage.connType {
	case ByDB:
		rows, err := storage.dbConn.Query(context.Background(), "SELECT short, original FROM urls WHERE userid=$1", userID)
		for rows.Next() {
			var rowValue UserURLs
			err := rows.Scan(&rowValue.ShortURL, &rowValue.OriginalURL)
			if err != nil {
				return nil, err
			}
			foundValues = append(foundValues, rowValue)
		}
		if err != nil {
			return nil, err
		} else {
			return foundValues, nil
		}
	case ByFile:
		db, err := os.OpenFile(storage.storageFilePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		mutex.RLock()
		scanner := bufio.NewScanner(db)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), userID) {
				value, _, _ := strings.Cut(scanner.Text(), ":_:")
				short, original, _ := strings.Cut(value, ":-:")
				foundValues = append(foundValues, UserURLs{ShortURL: config.GetServerConfig().BaseURL + "/" + short, OriginalURL: original})
			}
		}
		mutex.RUnlock()
		if err = db.Close(); err != nil {
			panic(err)
		}
		if foundValues != nil {
			return foundValues, nil
		} else {
			return nil, errors.New("values not found")
		}
	case byMap:
		return nil, errors.New("unable to open database")
	}
	return nil, nil
}

func (storage *Storage) SaveValue(key string, value string, userID string) {
	switch storage.connType {
	case ByDB:
		_, err := storage.dbConn.Exec(context.Background(), "INSERT INTO urls(short, original, userid) VALUES ($1, $2, $3)", key, value, userID)
		if err != nil {
			panic(err)
		}
		return
	case ByFile:
		db, err := os.OpenFile(storage.storageFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		_, err = db.WriteString(key + ":-:" + value + ":_:" + userID + "\n")
		if err != nil {
			panic(err)
		}
		mutex.Unlock()
		if err = db.Close(); err != nil {
			panic(err)
		}
		return
	case byMap:
		mutex.Lock()
		storage.dbMap[key] = value
		mutex.Unlock()
		return
	}
}

func (storage *Storage) BatchSave(values map[string]string, userID string) error {
	if storage.connType != ByDB {
		for k, v := range values {
			storage.SaveValue(k, v, userID)
		}
		return nil
	}
	batch := &pgx.Batch{}
	for k, v := range values {
		batch.Queue("INSERT INTO urls(short, original, userid) VALUES ($1, $2, $3)", k, v, userID)
	}
	bs := storage.dbConn.SendBatch(context.Background(), batch)
	_, err := bs.Exec()
	return err
}
