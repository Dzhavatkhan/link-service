package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

var ErrUrlExists = errors.New("url already exists")
var UrlNotFound = errors.New("url not found")


func New(ctx context.Context) (*Storage, error) {
	const op = "storage.postgres.New"

	connString := os.Getenv("CONN_STRING")
	if connString == "" {
		return nil, fmt.Errorf("%s: CONN_STRING environment variable is empty", op)
	}
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connection: %w", op, err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		alias VARCHAR(250) NOT NULL UNIQUE,
		url VARCHAR(250) NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create table: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.saveurl"
	var id int64
	query := `INSERT INTO urls(url,alias) VALUES ($1,$2) RETURNING id`
	err := s.pool.QueryRow(ctx, query, urlToSave, alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}
	return id, nil;

}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.geturl"
	var url string
	query := `SELECT url FROM urls WHERE alias=$1`
	err := s.pool.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("%s: failed to ping database: %w", op, err)
	}
	return url, nil;
}

func (s *Storage) DeleteURL(ctx context.Context, id int64) error {
	const op = "storage.postgres.deleteurl"
	query := `DELETE FROM urls WHERE id=$1`
	_,err := s.pool.Exec(ctx, query, id)
	if err != nil{
		return fmt.Errorf("%s: failed to ping database: %w", op, err)
	}
	return nil;
}

func (s *Storage) UpdateURL(ctx context.Context, id int64, urlToUpd string, aliasToUpd string) error{
	const op = "storage.postgres.UpdateURL"
	sqlQuery := `
	UPDATE urls
	SET url = $2, alias = $3
	WHERE id = $1
	`;
	_, err := s.pool.Exec(ctx, sqlQuery, id, urlToUpd, aliasToUpd);
	if err != nil {
		return fmt.Errorf("%s: failed to ping database: %w", op, err)
	}


	return nil;


}


func (s *Storage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}
