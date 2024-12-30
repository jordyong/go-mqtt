package database

import (
	"database/sql"
	"log"
	"os"
)

func openDB() {
	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer db.Close()
}

func storeMessage() {

}

func fetchMessages() {

}
