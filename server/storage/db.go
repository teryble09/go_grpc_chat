package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type dbPostgre struct {
	conn *sql.DB
}

func NewPostgresDBConnection(port string, dbname string, user string, password string) (dbPostgre, error) {
	cnnstring := fmt.Sprintf("port=%s dbname=%s user=%s password=%s", port, dbname, user, password)
	db, err := sql.Open("postgres", cnnstring)
	if err != nil {
		return dbPostgre{db}, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users
		(
			Username VARCHAR(10) PRIMARY KEY,
			HashPassword VARCHAR(60)
		)
	`)
	if err != nil {
		return dbPostgre{db}, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages
		(
			Id BIGSERIAL PRIMARY KEY,
			Sender VARCHAR(10),
			Content TEXT
		)
	`)

	return dbPostgre{db}, err
}
