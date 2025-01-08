package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	var err error
	dsn := "host=localhost port=5432 user=admin password=password dbname=taskdb sslmode=disable"
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to the database successfully!")
}
