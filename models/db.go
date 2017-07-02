package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite3 driver
)

// DB is the connection to gorm.
// Required to be exported so main.go can defer close the database.
var DB *gorm.DB

// InitDB will initialize the connection.
// It is used in main.go.
func InitDB() {
	var err error
	DB, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		panic("failed to connect database")
	}

	DB.Exec(`PRAGMA foreign_keys=ON`) // Need this to use foreign keys on sqlite.
	// defer db.Close()

	DB.LogMode(false)

	// Creates tables, columns, indexes. Does not delete or modify existing.
	DB.AutoMigrate(&Project{}, &Task{}, &Entry{})
}
