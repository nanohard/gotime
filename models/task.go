package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Task model.
type Task struct {
	gorm.Model
	Name      string
	Entries   []Entry
	Project   Project
	ProjectID uint `gorm:"index"`
	TotalTime time.Duration
}

// AllTasks queries the database for, and returns, all tasks after scanning them into a slice.
func AllTasks() []Task {
	t := []Task{}
	DB.Find(&t)
	return t
}

// GetTask queries the database for, and returns, one task
// after scanning it into the struct.
func GetTask(n string) Task {
	t := Task{}
	DB.Where("name = ?", n).First(&t)
	return t
}

// AddTask queries the database for one project by name.
// If the record exists then it is returned;
// else, it will create the record and return that one.
func AddTask(n string) Task {
	t := Task{Name: n}
	DB.FirstOrCreate(&t, Task{Name: n})
	return t
}
