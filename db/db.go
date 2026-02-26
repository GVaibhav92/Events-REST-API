package db

import (
	"REST-API/config"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB //gloabl db instance (thread-safe connection pool manager)

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", config.App.DBPath)

	if err != nil {
		log.Fatal("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10) // Maximum simultaneous database connections
	DB.SetMaxIdleConns(5)  // idle connections ready for reuse

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		panic("Couldn't enable foreign keys: " + err.Error())
	}
	createTables()
}

func createTables() {
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);
	`
	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER REFERENCES events(id),
		user_id INTEGER REFERENCES users(id)
	);
	`
	createRefreshTokensTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token TEXT NOT NULL UNIQUE,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err := DB.Exec(createEventsTable)
	if err != nil {
		panic("Could not create events table: " + err.Error())
	}

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		panic("Could not create users table: " + err.Error())
	}
	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		panic("Could not create registrations table: " + err.Error())
	}
	_, err = DB.Exec(createRefreshTokensTable)
	if err != nil {
		panic("Could not create refresh tokens table: " + err.Error())
	}
}
