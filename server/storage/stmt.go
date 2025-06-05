package storage

import "database/sql"

func InitStatemets(db *sql.DB) error {
	var err error
	stmtRegisterUser, err = db.Prepare(`
		INSERT INTO users (Username, HashPassword)
		VALUES ($0, $1)
	`)

	if err != nil {
		return err
	}

	stmtGetPassword, err = db.Prepare(`
		SELECT HashPassword FROM users
		WHERE Username = $1
	`)
	if err != nil {
		return err
	}

	stmtSaveMessage, err = db.Prepare(`
		INSERT INTO messages (Sender, Content)
		VALUES ($0, $1)
	`)
	if err != nil {
		return err
	}

	stmtGetMessages, err = db.Prepare(`
		SELECT Id, Sender, Content FROM messages
		WHERE Id >= $0 AND Id < $1
		ORDER BY Id ASC
	`)
	if err != nil {
		return err
	}

	stmtGetLastMessage, err = db.Prepare(`
		SELECT Id, Sender, Content FROM messages
		ORDER BY Id DESC LIMIT 1
	`)
	if err != nil {
		return err
	}

	return nil
}
