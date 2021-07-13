package db

import (
	"context"
	"github.com/jackc/pgx"
)

func NewDbConnector(ctx context.Context, login, pass, host, name string, port uint16) (*pgx.Conn, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     host,
		Port:     port,
		Database: name,
		Password: pass,
		User:     login,
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
