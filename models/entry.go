package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Entry model.
type Entry struct {
	gorm.Model
	Name      string
	Task      Task
	TaskID    uint `gorm:"index"`
	Project   Project
	ProjectID uint `gorm:"index"`
	TotalTime time.Duration
}

// AllEntries queries the database for, and returns, all entries after scanning them into a slice.
func AllEntries() []Entry {
	e := []Entry{}
	DB.Find(&e)
	return e
}

// GetEntry queries the database for, and returns, one entry
// after scanning it into the struct.
func GetEntry(n string) Entry {
	e := Entry{}
	DB.Where("name = ?", n).First(&e)
	return e
}

// AddEntry queries the database for one entry by name.
// If the record exists then it is returned;
// else, it will create the record and return that one.
func AddEntry(n string) Entry {
	e := Entry{Name: n}
	DB.FirstOrCreate(&e, Entry{Name: n})
	return e
}
