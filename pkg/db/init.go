package db

import (
	"database/sql"
	"fmt"
	"log"

	"chatbot/pkg/config"
	_ "github.com/go-sql-driver/mysql"
)

// InitDB initializes and returns the database connection.
func InitDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	log.Println("Connected to the database successfully.")
	return db
}
