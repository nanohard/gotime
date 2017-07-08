package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// CurrentTask holds the struct of the currently
// highlighted task.
var CurrentTask Task

// Task model.
type Task struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Entries     []Entry
	Project     Project
	ProjectID   uint `gorm:"index;not null"`
	TotalTime   time.Duration
}

// AllTasks queries the database for, and returns, all tasks after scanning them into a slice.
func AllTasks(p Project) []Task {
	var t []Task
	DB.Model(&p).Related(&t)
	return t
}

// GetTask queries the database for, and returns, one task
// after scanning it into the struct.
func GetTask(n string) Task {
	var t Task
	n = strings.TrimSpace(n)
	// DB.Where(&Task{Name: n}).First(&t)
	DB.Model(&CurrentProject).Where(&Task{Name: n}).Related(&t)
	return t
}

// AddTask queries the database for one project by name.
// If the record exists then it is returned;
// else, it will create the record and return that one.
func AddTask(n string, p Project) Task {
	var t Task
	n = strings.TrimSpace(n)
	DB.FirstOrCreate(&t, Task{Name: n, ProjectID: p.ID})
	return t
}

// Delete one task and its children.
func (t Task) Delete() {
	entries := AllEntries(t)
	for _, e := range entries {
		e.Delete()
	}
	DB.Delete(&t)
}

// HoursMinutes takes all entries for a task and returns its
// hours and minutes, both as ints.
func (t Task) HoursMinutes() (h int, m int) {
	e := AllEntries(t)
	var total float64
	for _, i := range e {
		total += i.TotalTime.Seconds()
	}
	hours := int(total) / 3600
	f := float64((int(total) % 3600.0) / 60.0)
	i := float64(f) + float64(0.5)
	minutes := int(i)
	return hours, minutes
}
