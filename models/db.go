package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite3 driver
)

// DB is the connection to gorm.
// Required to be exported so main.go can defer close the database.
var DB *gorm.DB

// Setting will contain the only row from the settings table.
type setting struct {
	gorm.Model
	SortBy    string
	SortOrder string
}

// Setting will contain the only row from the settings table.
var Setting setting

// InitDB will initialize the connection.
// It is used in main.go.
func InitDB() {
	var err error
	DB, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		panic("failed to connect database")
	}

	DB.Exec(`PRAGMA foreign_keys=ON`) // Need this to use foreign keys on sqlite.
	// DB.LogMode(true) // Turn on for debugging.

	// Creates tables, columns, indexes. Does not delete or modify existing.
	DB.AutoMigrate(&Project{}, &Task{}, &Entry{}, &setting{})

	// Create settings table.
	DB.Exec("INSERT OR IGNORE INTO settings (id, sort_by, sort_order) VALUES(1, 'name', 'asc')")

	// Get only row from settings table and insert into exported variable.
	row := DB.Table("settings").Where("id = ?", "1").Select("sort_by, sort_order").Row() // (*sql.Row)
	row.Scan(&Setting.SortBy, &Setting.SortOrder)
	// DB.First does not work for some reason...
	// DB.First(&S, 1)
}
