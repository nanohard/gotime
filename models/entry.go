package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var CurrentEntry Entry

// Entry model.
type Entry struct {
	gorm.Model
	Name      time.Time
	Task      Task
	TaskID    uint `gorm:"index;not null"`
	Project   Project
	ProjectID uint `gorm:"index;not null"`
	TotalTime time.Duration
}

// AllEntries queries the database for, and returns, all entries after scanning them into a slice.
func AllEntries() []Entry {
	var e []Entry
	DB.Find(&e)
	return e
}

// GetEntry queries the database for, and returns, one entry
// after scanning it into the struct.
func GetEntry(n string) Entry {
	var e Entry
	n = strings.TrimSpace(n)
	DB.Where("name = ?", n).First(&e)
	return e
}

// AddEntry queries the database for one entry by name.
// If the record exists then it is returned;
// else, it will create the record and return that one.
func AddEntry() Entry {
	t := time.Now()
	var e Entry
	DB.FirstOrCreate(&e, Entry{Name: t})
	return e
}
