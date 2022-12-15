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

func (file *File) GetValue(_ context.Context, key string) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	db, err := os.OpenFile(file.Path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	foundValue := ""
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), key) {
			_, foundValue, _ = strings.Cut(scanner.Text(), ":-:")
			foundValue, _, _ = strings.Cut(foundValue, ":_:")
			break
		}
	}
	if err = db.Close(); err != nil {
		return foundValue, err
	}
	if foundValue != "" {
		return foundValue, nil
	}
	return "", NotFoundError()
}

func (file *File) GetValuesByID(_ context.Context, userID string) ([]UserURL, error) {
	foundValues := make([]UserURL, 0, 16)
	mutex.RLock()
	defer mutex.RUnlock()
	db, err := os.OpenFile(file.Path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), userID) {
			value, _, _ := strings.Cut(scanner.Text(), ":_:")
			short, original, _ := strings.Cut(value, ":-:")
			foundValues = append(foundValues, UserURL{ShortURL: config.GetServerConfig().BaseURL + "/" + short, OriginalURL: original})
		}
	}
	if err = db.Close(); err != nil {
		return foundValues, err
	}
	if len(foundValues) != 0 {
		return foundValues, nil
	}
	return nil, NotFoundError()
}

func (file *File) SaveValue(_ context.Context, key string, value string, userID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	db, err := os.OpenFile(file.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	_, err = db.WriteString(key + ":-:" + value + ":_:" + userID + "\n")
	if err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}
	return nil
}

func (file *File) BatchSave(_ context.Context, values map[string]string, userID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	db, err := os.OpenFile(file.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	for k, v := range values {
		_, err = db.WriteString(k + ":-:" + v + ":_:" + userID + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (file *File) FindDuplicate(_ context.Context, value string) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	db, err := os.OpenFile(file.Path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	key := ""
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), value) {
			key, _, _ = strings.Cut(scanner.Text(), ":-:")
			break
		}
	}
	if err = db.Close(); err != nil {
		return "", err
	}
	if key != "" {
		return key, FoundDuplicateError()
	}
	return "", NotFoundError()
}
