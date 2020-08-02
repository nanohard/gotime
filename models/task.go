package models

import (
	"encoding/csv"
	"github.com/jroimartin/gocui"
	"log"
	"os"
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

// ExportTaskCsv creates a CSV file of all of the user's entries for the current task. Bound to Ctrl-E.
func ExportTaskCsv(g *gocui.Gui, v *gocui.View) error {
	var err error = nil

	// get all entries from the database
	entries := AllEntries(CurrentTask)

	// create CSV file
	taskFileName := strings.ReplaceAll(strings.ToLower(CurrentTask.Name), " ", "_")
	file, err := os.Create(taskFileName + "_entries.csv")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// create a new writer type
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write CSV headers to file
	headers := []string{"Entry Details", "Time Worked", "Project Name", "Task Name"}
	err = writer.Write(headers)

	// find associated project for the task from the database
	var p Project
	DB.First(&p, CurrentTask.ProjectID)

	// write contents of each entry to a row in the CSV (string-writeable data)
	for _, record := range entries {
		timeWorked := record.End.Sub(record.Start)

		values := []string{record.Details, timeWorked.String(), p.Name, CurrentTask.Name}
		err := writer.Write(values)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}
