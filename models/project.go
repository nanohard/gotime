package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// CurrentProject holds the struct of the currently
// highlighted project.
var CurrentProject Project

// Project model.
type Project struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);unique;not null"`
	Description string
	Tasks       []Task
	Entries     []Entry
	TotalTime   time.Duration
}

// AllProjects queries the database for, and returns, all projects
// after scanning them into a slice of structs.
func AllProjects() []Project {
	var p []Project
	o := Setting.SortBy + " " + Setting.SortOrder
	DB.Order(o).Find(&p)
	return p
}

// GetProject queries the database for, and returns, one project
// after scanning it into the struct.
func GetProject(n string) Project {
	var p Project
	n = strings.TrimSpace(n)
	DB.Where(&Project{Name: n}).First(&p)
	return p
}

// AddProject queries the database for one project by name (project names are unique).
// If the record exists then it is returned;
// else, it will create the record and return that one.
func AddProject(n string) Project {
	var p Project
	n = strings.TrimSpace(n)
	DB.FirstOrCreate(&p, Project{Name: n})
	return p
}

// Delete one project and its children.
func (p Project) Delete() {
	tasks := AllTasks(p)
	for _, t := range tasks {
		t.Delete()
	}
	DB.Delete(&p)
}

// UpdateProject will update the project with values defined outside.
func UpdateProject(p Project) Project {
	DB.Save(&p)
	return p
}

// HoursMinutes takes all entries for a project and returns its
// hours and minutes, both as ints.
func (p Project) HoursMinutes() (h int, m int) {
	// var e []Entry
	// DB.Model(&p).Related(&e)
	DB.Preload("Entries").Find(&p)
	var total float64
	for _, i := range p.Entries {
		total += i.TotalTime.Seconds()
	}
	hours := int(total) / 3600
	f := float64((int(total) % 3600.0) / 60.0)
	i := float64(f) + float64(0.5)
	minutes := int(i)
	return hours, minutes
}
