package util

import (
	"log"
	"os"
	"database/sql"
    "path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func connectDB() {
	migrateDB()

	testDB, err := sql.Open("sqlite3", string(os.Getenv("DB_PATH")))
	if err != nil {
		log.Panic(err)
	}

	db = testDB
}

func migrateDB() {
	testDB, err := sql.Open("sqlite3", string(os.Getenv("DB_PATH")))
	if err != nil {
		log.Panic(err)
	}
	db = testDB
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS devices (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            mac_address VARCHAR(64) NOT NULL,
            blacklisted BOOLEAN NOT NULL CHECK (blacklisted in (0, 1)),
            UNIQUE(mac_address)
        );
    `)
	if err != nil {
		log.Panic(err)
	}
}

func createDB() {
	path := os.Getenv("DB_PATH")
    os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(os.Getenv("DB_PATH"))
		log.Printf("[DB]: Created database file as %s", path)
	}
}

func init() {
	godotenv.Load(".local.env")
	createDB()
	connectDB()
	log.Println("[DB]: Connected to Database")
}
