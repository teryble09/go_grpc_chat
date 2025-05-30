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
	cnnstring := fmt.Sprintf("port=%s dbname=%s user=%s password=%s sslmode=disable", port, dbname, user, password)
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

	stmtRegisterUser, err = db.Prepare(`
		INSERT INTO users (Username, HashPassword)
		VALUES ($0, $1)
	`)

	if err != nil {
		return dbPostgre{db}, err
	}

	stmtGetPassword, err = db.Prepare(`
		SELECT HashPassword FROM users
		WHERE Username = $1
	`)
	if err != nil {
		return dbPostgre{db}, err
	}

	stmtSaveMessage, err = db.Prepare(`
		INSERT INTO messages (Sender, Content)
		VALUES ($0, $1)
	`)
	if err != nil {
		return dbPostgre{db}, err
	}

	stmtGetMessages, err = db.Prepare(`
		SELECT Id, Sender, Content FROM messages
		WHERE Id >= $0 AND Id < $1
		ORDER BY Id ASC
	`)
	if err != nil {
		return dbPostgre{db}, err
	}

	stmtGetLastMessage, err = db.Prepare(`
		SELECT Id, Sender, Content FROM messages
		ORDER BY Id DESC LIMIT 1

	`)
	if err != nil {
		return dbPostgre{db}, err
	}

	return dbPostgre{db}, err
}
