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
	Name      string `gorm:"type:varchar(100);unique;not null"`
	Tasks     []Task
	Entries   []Entry
	TotalTime time.Duration
}

// AllProjects queries the database for, and returns, all projects
// after scanning them into a slice of structs.
func AllProjects() []Project {
	var p []Project
	o := Setting.Sort + " " + Setting.SortOrder
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
