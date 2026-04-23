// Пакет storage реализует слой хранения данных системы Gophermart на базе PostgreSQL.
package storage

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrNotFound = errors.New("object not found")

// PostgresStorage реализует все интерфейсы хранилища поверх соединения с PostgreSQL.
type PostgresStorage struct {
	db  *sql.DB
	dsn string
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db, dsn: dsn}, nil
}

func (p *PostgresStorage) Migrate(migrationsPath string) error {
	m, err := migrate.New(migrationsPath, p.dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
