package storage

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type DB struct {
	dbPool *pgxpool.Pool
}

func NewDB() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pgPool, err := pgxpool.New(ctx, config.GetServerConfig().DBAddress)
	if err != nil {
		log.Fatal("Unable to create DB: ", err)
	}
	return &DB{dbPool: pgPool}
}

func (db *DB) GetValue(ctx context.Context, key string) (string, error) {
	dbConn, err := db.dbPool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer dbConn.Release()
	var value string
	err = dbConn.QueryRow(ctx, "SELECT original FROM urls WHERE short=$1", key).Scan(&value)
	if err != nil {
		return "", NotFoundError()
	}
	return value, nil
}

func (db *DB) GetValuesByID(ctx context.Context, userID string) ([]UserURL, error) {
	foundValues := make([]UserURL, 0, 16)

	dbConn, err := db.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer dbConn.Release()
	rows, err := dbConn.Query(ctx, "SELECT short, original FROM urls WHERE userid=$1", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowValue UserURL
		err := rows.Scan(&rowValue.ShortURL, &rowValue.OriginalURL)
		if err != nil {
			return nil, err
		}
		foundValues = append(foundValues, rowValue)
	}
	if len(foundValues) != 0 {
		return foundValues, nil
	}
	return foundValues, nil
}

func (db *DB) SaveValue(ctx context.Context, key string, value string, userID string) error {
	dbConn, err := db.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer dbConn.Release()
	_, err = dbConn.Exec(ctx, "INSERT INTO urls(short, original, userid) VALUES ($1, $2, $3)", key, value, userID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) BatchSave(ctx context.Context, values map[string]string, userID string) error {

	batch := &pgx.Batch{}
	for k, v := range values {
		batch.Queue("INSERT INTO urls(short, original, userid) VALUES ($1, $2, $3)", k, v, userID)
	}
	dbConn, err := db.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer dbConn.Release()
	bs := dbConn.SendBatch(ctx, batch)
	_, err = bs.Exec()
	return err
}

func (db *DB) FindDuplicate(ctx context.Context, value string) (string, error) {
	dbConn, err := db.dbPool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer dbConn.Release()
	var key string
	err = dbConn.QueryRow(ctx, "SELECT short FROM urls WHERE original=$1", value).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, FoundDuplicateError()
}
