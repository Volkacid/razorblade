package storage

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/jackc/pgx/v5"
	"time"
)

func CheckDBConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, config.GetServerConfig().DBAddress)
	if err == nil {
		defer conn.Close(ctx)
		err = conn.Ping(ctx)
	}
	return err == nil
}

func InitializeDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, config.GetServerConfig().DBAddress)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS urls (short varchar(30), original varchar(300), userid varchar(100), PRIMARY KEY (short))")
	if err != nil {
		return err
	}
	return nil
}
