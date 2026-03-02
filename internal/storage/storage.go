package storage

import (
	"billohub/pkg/helper"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Storage implements the model.AgentStorage interface for PostgreSQL using GORM.
type Storage struct {
	db *gorm.DB
}

// NewStorage creates a new Storage instance with GORM.
func NewStorage(dsn string) (*Storage, error) {
	var err error
	var db *gorm.DB
	if strings.HasPrefix(dsn, "postgres") {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		if dsn == "" {
			dsn = "my.db"
		}
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, helper.WrapError(err, "failed to connect to database with gorm")
	}

	return &Storage{db: db}, nil
}

// DB returns the underlying gorm.DB instance.
func (s *Storage) DB() *gorm.DB {
	return s.db
}
