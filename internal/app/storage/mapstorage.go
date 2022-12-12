package storage

import "context"

type IDValue struct {
	OrigURL string
	UserID  string
}

type dbMap struct {
	db map[string]IDValue
}

func NewMapDB() *dbMap {
	return &dbMap{db: make(map[string]IDValue)}
}

func (db *dbMap) GetValue(key string) (string, error) {
	mutex.RLock()
	value, ok := db.db[key]
	mutex.RUnlock()
	if !ok {
		return "", NotFoundError()
	}
	return value.OrigURL, nil
}

func (db *dbMap) GetValuesByID(userID string) ([]UserURL, error) {
	var foundValues []UserURL

	for k, v := range db.db {
		if v.UserID == userID {
			foundValues = append(foundValues, UserURL{OriginalURL: v.OrigURL, ShortURL: k})
		}
	}
	if foundValues != nil {
		return foundValues, nil
	}
	return nil, NotFoundError()
}

func (db *dbMap) SaveValue(key string, value string, userID string) error {
	mutex.Lock()
	db.db[key] = IDValue{OrigURL: value, UserID: userID}
	mutex.Unlock()
	return nil
}

func (db *dbMap) BatchSave(ctx context.Context, values map[string]string, userID string) error {
	for k, v := range values {
		err := db.SaveValue(k, v, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *dbMap) FindDuplicate(value string) (string, error) {
	mutex.RLock()
	for k, v := range db.db {
		if v.OrigURL == value {
			return k, nil
		}
	}
	mutex.RUnlock()
	return "", NotFoundError()
}
