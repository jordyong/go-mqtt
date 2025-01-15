package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteService struct {
	db *sql.DB
}

func NewSQliteService() (*SQLiteService, error) {
	db, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	return &SQLiteService{db: db}, nil
}

func (dbService *SQLiteService) testConnection() error {
	for i := 0; i < 3; i++ { // Retry logic
		if err := dbService.db.Ping(); err == nil {
			log.Println("connected to database")
			return nil
		}
		log.Printf("failed to ping database"+" (attempt %d)", i+1)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	return fmt.Errorf("failed to ping database after multiple attempts")
}

func (dbService *SQLiteService) GetDB() *sql.DB {
	dbService.testConnection()
	return dbService.db
}

func (dbService *SQLiteService) Disconnect() {
	if dbService.db != nil {
		dbService.db.Close()
		log.Println("Database connection closed")
	}
}

func (dbService *SQLiteService) CreateTable() {
	createDevicesTable := `CREATE TABLE devices (
    device_id TEXT PRIMARY KEY,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	  );`

	createDataTable := `CREATE TABLE data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    data_type TEXT,
    data_value REAL,
    FOREIGN KEY (device_id) REFERENCES devices(device_id)
  );`

	db := dbService.GetDB()
	statement, err := db.Prepare(createDevicesTable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements

	statement, err = db.Prepare(createDataTable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

// We are passing db reference connection from main to our method with other parameters
func (dbService *SQLiteService) InsertStudent(
	code string,
	name string,
	program string,
) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	db := dbService.GetDB()
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (dbService *SQLiteService) DisplayDevice() {
	db := dbService.GetDB()
	row, err := db.Query("SELECT * FROM devices ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var status string
		row.Scan(&id, &status)
		log.Println("Device: ", id, " ", status)
	}
}
