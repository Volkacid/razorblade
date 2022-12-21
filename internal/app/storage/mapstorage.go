package storage

import (
	"context"
)

type IDValue struct {
	OrigURL   string
	UserID    string
	IsDeleted bool
}

type dbMap struct {
	db map[string]IDValue
}

func NewMapDB() *dbMap {
	return &dbMap{db: make(map[string]IDValue)}
}

func (db *dbMap) GetValue(_ context.Context, key string) (string, error) {
	mutex.RLock()
	value, ok := db.db[key]
	mutex.RUnlock()
	if !ok {
		return "", NotFoundError()
	}
	if value.IsDeleted {
		return "", ValueDeletedError()
	}
	return value.OrigURL, nil
}

func (db *dbMap) GetValuesByID(_ context.Context, userID string) ([]UserURL, error) {
	foundValues := make([]UserURL, 0, 16)
	mutex.RLock()
	defer mutex.RUnlock()
	for k, v := range db.db {
		if v.UserID == userID && !v.IsDeleted {
			foundValues = append(foundValues, UserURL{OriginalURL: v.OrigURL, ShortURL: k})
		}
	}
	if len(foundValues) != 0 { //Necessary for correct http 204 status handling
		return foundValues, nil
	}
	return nil, NotFoundError()
}

func (db *dbMap) SaveValue(_ context.Context, key string, value string, userID string) error {
	mutex.Lock()
	db.db[key] = IDValue{OrigURL: value, UserID: userID, IsDeleted: false}
	mutex.Unlock()
	return nil
}

func (db *dbMap) BatchSave(_ context.Context, values map[string]string, userID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for k, v := range values {
		db.db[k] = IDValue{OrigURL: v, UserID: userID, IsDeleted: false}
	}
	return nil
}

func (db *dbMap) FindDuplicate(_ context.Context, value string) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	for k, v := range db.db {
		if v.OrigURL == value {
			return k, nil
		}
	}
	return "", NotFoundError()
}

func (db *dbMap) DeleteURLs(urls []string) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, key := range urls {
		userVal := db.db[key]
		userVal.IsDeleted = true
		db.db[key] = userVal
	}
}
