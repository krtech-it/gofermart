// Пакет storage реализует слой хранения данных системы Gophermart на базе PostgreSQL.
package storage

import (
	"database/sql"
	"errors"
)

// PostgresStorage реализует все интерфейсы хранилища поверх соединения с PostgreSQL.
type PostgresStorage struct {
	db *sql.DB
}

var ErrNotFound = errors.New("user not found")
