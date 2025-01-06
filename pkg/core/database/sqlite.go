package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteService struct {
	db *sql.DB
}

func InitDB() (*sql.DB, error) {
	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		return nil, err
	}
	file.Close()
	log.Println("sqlite-database.db created")

	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer db.Close()

	// SQL statement to create the todos table if it doesn't exist
	sqlStmt := `
 CREATE TABLE IF NOT EXISTS test (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  title TEXT
 );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf(
			"Error creating table: %q: %s\n",
			err,
			sqlStmt,
		) // Log an error if table creation fails
	}

	return db, nil
}

func StoreMessage() {

}

func FetchMessages() {

}
