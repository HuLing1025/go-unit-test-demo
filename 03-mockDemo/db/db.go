package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

const DB_DSN = "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"

// dbClient private static variable.
var dbClient *gorm.DB

func init() {
	db, err := newDbClient()
	if err != nil {
		fmt.Printf("connect to db error: %v", err)
	}

	dbClient = db
}

func GetDB() (db *gorm.DB, err error) {
	if dbClient != nil {
		return dbClient, nil
	}

	db, err = newDbClient()
	if err != nil {
		fmt.Printf("connect to db error: %v", err)
	}

	dbClient = db

	return
}

func newDbClient() (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Error),
	})
	return
}
