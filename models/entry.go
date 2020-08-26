package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// TL is the time layout we will use.
	TL = "2006-01-02 15:04:05"
)

// CurrentEntry holds the struct of the currently
// highlighted entry.
var CurrentEntry Entry

// Entry model.
type Entry struct {
	gorm.Model
	Name      string `gorm:"index;unique;not null"`
	Details   string
	Start     time.Time
	End       time.Time
	Task      Task
	TaskID    uint `gorm:"index;not null"`
	Project   Project
	ProjectID uint `gorm:"index;not null"`
	TotalTime time.Duration
}

// AllEntries queries the database for, and returns, all entries after scanning them into a slice.
func AllEntries(t Task) []Entry {
	var e []Entry
	// SoryBy = "name" and timestamps are used as the names for entries,
	// so by sorting in desc order it will put the latest entry at the top.
	o := Setting.SortBy + " " + "desc"
	DB.Order(o).Model(&t).Related(&e)
	return e
}

// GetEntry queries the database for, and returns, one entry
// after scanning it into the struct.
func GetEntry(n string) Entry {
	var e Entry
	DB.Where(&Entry{Name: n}).First(&e)
	return e
}

// Delete one entry.
func (e Entry) Delete() {
	DB.Delete(&e)
}

// StartEntry queries the database for one entry by name.
// If the record exists then it is returned;
// else, it will create the record and return that one.
// It takes a task and start time as arguments.
func StartEntry(t Task, s time.Time) Entry {
	n := TimeIn(s)
	var e Entry
	DB.FirstOrCreate(&e, Entry{Name: n, Start: s, TaskID: t.ID, ProjectID: CurrentProject.ID})
	return e
}

// StopEntry updates the entry with the stopped time, duration, and details.
// It takes an entry struct and stop time as arguments.
func StopEntry(e Entry, s time.Time, d string) Entry {
	var tt time.Duration = s.Sub(e.Start)
	DB.Model(&e).Updates(Entry{End: s, TotalTime: tt, Details: d})
	return e
}

// HoursMinutes takes an Entry{} and returns its
// hours and minutes, both as ints.
func (e Entry) HoursMinutes() (h int, m int) {
	total := int(e.TotalTime.Seconds())
	hours := total / 3600
	f := float64((total % 3600.0) / 60.0)
	i := float64(f) + float64(0.5)
	minutes := int(i)
	return hours, minutes
}

// TimeIn turns a time object into a datetime string.
func TimeIn(t time.Time) string {
	return t.Format(TL)
}

// TimeOut turns a datetime string into a time object.
func TimeOut(s string) time.Time {
	t, _ := time.Parse(TL, s)
	return t
}
