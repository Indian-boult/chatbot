package db

import (
	"chatbot/pkg/model"
	"database/sql"
	"fmt"
)

// Database represents the database operations interface.
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new Database instance.
func NewDatabase(db *sql.DB) *Database {
	return &Database{DB: db}
}

// SaveImage stores an image's metadata in the database.
func (db *Database) SaveImage(image *model.ImageInfo) error {
	query := `INSERT INTO images (identifier, url) VALUES (?, ?)`
	_, err := db.DB.Exec(query, image.Identifier, image.URL)
	if err != nil {
		return fmt.Errorf("SaveImage: %v", err)
	}
	return nil
}

// GetImage retrieves an image's metadata using its identifier.
func (db *Database) GetImage(identifier string) (*model.ImageInfo, error) {
	query := `SELECT id, identifier, url FROM images WHERE identifier = ? LIMIT 1`
	row := db.DB.QueryRow(query, identifier)

	var image model.ImageInfo
	if err := row.Scan(&image.ID, &image.Identifier, &image.URL); err != nil {
		return nil, fmt.Errorf("GetImage: %v", err)
	}
	return &image, nil
}
