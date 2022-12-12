package storage

import (
	"bufio"
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"os"
	"strings"
)

type File struct {
	Path string
}

func NewFileDB(filePath string) *File {
	return &File{Path: filePath}
}

func (file *File) GetValue(key string) (string, error) {
	db, err := os.OpenFile(file.Path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
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
		return foundValue, err
	}
	if foundValue != "" {
		return foundValue, nil
	}
	return "", NotFoundError()
}

func (file *File) GetValuesByID(userID string) ([]UserURL, error) {
	var foundValues []UserURL

	db, err := os.OpenFile(file.Path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	mutex.RLock()
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), userID) {
			value, _, _ := strings.Cut(scanner.Text(), ":_:")
			short, original, _ := strings.Cut(value, ":-:")
			foundValues = append(foundValues, UserURL{ShortURL: config.GetServerConfig().BaseURL + "/" + short, OriginalURL: original})
		}
	}
	mutex.RUnlock()
	if err = db.Close(); err != nil {
		return foundValues, err
	}
	if foundValues != nil {
		return foundValues, nil
	}
	return nil, NotFoundError()
}

func (file *File) SaveValue(key string, value string, userID string) error {
	db, err := os.OpenFile(file.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	mutex.Lock()
	_, err = db.WriteString(key + ":-:" + value + ":_:" + userID + "\n")
	if err != nil {
		return err
	}
	mutex.Unlock()
	if err = db.Close(); err != nil {
		return err
	}
	return nil
}

func (file *File) BatchSave(ctx context.Context, values map[string]string, userID string) error {
	for k, v := range values {
		err := file.SaveValue(k, v, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (file *File) FindDuplicate(value string) (string, error) {
	db, err := os.OpenFile(file.Path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	key := ""
	mutex.RLock()
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), value) {
			key, _, _ = strings.Cut(scanner.Text(), ":-:")
			break
		}
	}
	mutex.RUnlock()
	if err = db.Close(); err != nil {
		return "", err
	}
	if key != "" {
		return key, FoundDuplicateError()
	}
	return "", NotFoundError()
}
