package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(host string, port uint, user string, password string) (*Postgres, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable database=file_service", host, port, user, password))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
