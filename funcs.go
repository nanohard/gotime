package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/nanohard/gotime/models"
	"github.com/pkg/errors"
)

var bottom = false
var dView *gocui.View

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	// Check to make sure data exists in the next line,
	// otherwise disallow scroll down.
	if v != nil && lineBelow(g, v) == true {
		v.MoveCursor(0, 1, false)
		_, cy := v.Cursor()
		n, _ := v.Line(cy)
		if v.Name() == P {
			nv, _ := g.View(T)
			if n != "" {
				models.CurrentProject = models.GetProject(n)
			}
			// log.Println("cursorDown project.id:", models.CurrentProject.ID)
			redrawTasks(g, nv)
		} else if v.Name() == T {
			nv, _ := g.View(E)
			if n != "" {
				models.CurrentTask = models.GetTask(n)
			}
			// log.Println("cursorDown task.id:", models.CurrentTask.ID)
			redrawEntries(g, nv)
		} else if v.Name() == E {
			nv, _ := g.View(O)
			if n != "" {
				models.CurrentEntry = models.GetEntry(n)
			}
			// log.Println("cursorUp entry.id:", models.CurrentEntry.ID)
			redrawOutput(g, nv)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, -1, false)
		_, cy := v.Cursor()
		n, _ := v.Line(cy)
		if v.Name() == P {
			nv, _ := g.View(T)
			if n != "" {
				models.CurrentProject = models.GetProject(n)
			} else {
				models.CurrentProject = models.Project{}
			}
			// log.Println("cursorUp project.id:", models.CurrentProject.ID)
			redrawTasks(g, nv)
		} else if v.Name() == T {
			nv, _ := g.View(E)
			if n != "" {
				models.CurrentTask = models.GetTask(n)
			} else {
				models.CurrentTask = models.Task{}
			}
			// log.Println("cursorUp task.id:", models.CurrentTask.ID)
			redrawEntries(g, nv)
		} else if v.Name() == E {
			nv, _ := g.View(O)
			if n != "" {
				models.CurrentEntry = models.GetEntry(n)
			} else {
				models.CurrentEntry = models.Entry{}
			}
			// log.Println("cursorUp entry.id:", models.CurrentEntry.ID)
			redrawOutput(g, nv)
		}
	}
	return nil
}

// Returns true if there is a non-empty string in cursor position y+1.
// Otherwise returns false.
func lineBelow(g *gocui.Gui, v *gocui.View) bool {
	_, cy := v.Cursor()
	if l, _ := v.Line(cy + 1); l != "" {
		return true
	}
	return false
}

// Copy the input view (iv) and handle it.
// Used to add project or task.
func copyInput(g *gocui.Gui, iv *gocui.View) error {
	var err error
	// We want to read the view’s buffer from the beginning.
	iv.Rewind()
	// Get the output view via its name.
	var ov *gocui.View
	// If there is text input then add the item,
	// else go back to the input view.
	switch iv.Name() {
	case "addProject":
		ov, _ = g.View(P)
		if iv.Buffer() != "" {
			models.AddProject(iv.Buffer())
		} else {
			inputView(g, ov)
			return nil
		}
	case "addTask":
		ov, _ = g.View(T)
		if iv.Buffer() != "" {
			models.AddTask(iv.Buffer(), models.CurrentProject)
		} else {
			inputView(g, ov)
			return nil
		}
	}
	// Clear the input view
	iv.Clear()
	// No input, no cursor.
	g.Cursor = false
	// !!!
	// Must delete keybindings before the view, or fatal error !!!
	// !!!
	g.DeleteKeybindings(iv.Name())
	if err = g.DeleteView(iv.Name()); err != nil {
		return err
	}
	// Set the view back.
	if _, err = g.SetCurrentView(ov.Name()); err != nil {
		return err
	}
	switch ov.Name() {
	case P:
		redrawProjects(g, ov)
	case T:
		redrawTasks(g, ov)
	}
	return err
}

// Add item to the current view (cv) using the text from the input view (iv).
func inputView(g *gocui.Gui, cv *gocui.View) error {
	maxX, maxY := g.Size()
	var title string
	var name string
	switch cv.Name() {
	case P:
		title = "Name of new project"
		name = "addProject"
	case T:
		title = "Name of new task"
		name = "addTask"
	}
	if iv, err := g.SetView(name, maxX/2-12, maxY/2, maxX/2+12, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		iv.Title = title
		iv.Editable = true
		g.Cursor = true
		if _, err := g.SetCurrentView(name); err != nil {
			return err
		}
		if err := g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, copyInput); err != nil {
			return err
		}
	}
	return nil
}

// Get the current view (cv) and transfer cursor to the new view (nv).
// Disallow if there is no string at current cursor.
func selectItem(g *gocui.Gui, cv *gocui.View) error {
	var err error
	_, cy := cv.Cursor()
	n, _ := cv.Line(cy)
	// If line at cursor is not empty (item is selected) then continue.
	if n != "" {
		var nv *gocui.View
		switch cv.Name() {
		case P:
			if nv, err = g.SetCurrentView(T); err != nil {
				return err
			}
			models.CurrentEntry = models.Entry{}
			// log.Println("selectItem project view CurrentEntry:", models.CurrentEntry.ID)
			nv.SetCursor(0, 0)
			cursorUp(g, nv)
		case T:
			if nv, err = g.SetCurrentView(E); err != nil {
				return err
			}
			models.CurrentEntry = models.Entry{}
			// log.Println("selectItem task view CurrentEntry:", models.CurrentEntry.ID)
			nv.SetCursor(0, 0)
			cursorUp(g, nv)
		}
		// Turn on highlight and set cursor to 0,0 of the new view.
		nv.Highlight = true
		if err = nv.SetCursor(0, 0); err != nil {
			return err
		}
	}
	return nil
}

// Get the name of the item at the cursor and delete it.
// Disallow if there is no string at current cursor.
// Hmm... I'm thinking that if two tasks have the same name it would randomly pick one.
// Maybe not, but I forget how it sets the models.Current[Item]. Something to look out for.
func deleteItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	n, _ := v.Line(cy)
	// If line at cursor is not empty (item is selected) then continue.
	if n != "" {
		switch v.Name() {
		case P:
			models.CurrentProject.Delete()
			models.CurrentProject = models.Project{}
			redrawProjects(g, v)
			v.SetCursor(0, 0)
			cursorUp(g, v)
		case T:
			models.CurrentTask.Delete()
			models.CurrentTask = models.Task{}
			redrawTasks(g, v)
			v.SetCursor(0, 0)
			cursorUp(g, v)

		case E:
			models.CurrentEntry.Delete()
			models.CurrentEntry = models.Entry{}
			redrawEntries(g, v)
			v.SetCursor(0, 0)
			cursorUp(g, v)
		}
	}
	return err
}

// Get the current view (cv) and transfer cursor to the new view (nv).
// Basically the opposite of selectItem.
func goBack(g *gocui.Gui, cv *gocui.View) error {
	var err error
	var nv *gocui.View
	switch cv.Name() {
	// Move from tasks to projects view.
	case T:
		if nv, err = g.SetCurrentView(P); err != nil {
			return err
		}
		entriesView, _ := g.View(E)
		redrawEntries(g, entriesView)
	// Move from entries to tasks view.
	case E:
		if nv, err = g.SetCurrentView(T); err != nil {
			return err
		}
		models.CurrentEntry = models.Entry{}
		outputView, _ := g.View(O)
		redrawOutput(g, outputView)
	}
	// Turn off highlight of current view and make sure it's on for the new view.
	cv.Highlight = false
	// Probably redundant.
	nv.Highlight = true
	return nil
}

// Get the view and redraw it with current database info.
// It's important to note that this function will call
// redrawTasks, which will call
// redrawEntries, which will call
// redrawOutput. Make it fucking rain.
func redrawProjects(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	items := models.AllProjects()
	// Loop through projects to add their names to the view.
	for _, i := range items {
		// We can simply Fprint to a view.
		_, err := fmt.Fprintln(v, i.Name)
		if err != nil {
			log.Println("Error writing to the projects view:", err)
		}
	}
	// If there are no projects then one must be created.
	if len(items) == 0 {
		inputView(g, v)
	}
	// While the text may shift lines on insert the cursor does not,
	// so we need to refresh the tasks view with the currently highlighted project.
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	if l != "" {
		models.CurrentProject = models.GetProject(l)
	}
	// log.Println("redrawProjects project.id:", models.CurrentProject.ID)
	tasksView, _ := g.View(T)
	// Projects is only redrawn if in the projects view, so it's
	// safe to zero the current task and entry.
	models.CurrentTask = models.Task{}
	models.CurrentEntry = models.Entry{}
	redrawTasks(g, tasksView)
	tasksView.Highlight = false
}

// Get the view and redraw it with current database info.
func redrawTasks(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	items := models.AllTasks(models.CurrentProject)
	// Loop through tasks to add their names to the view.
	for _, i := range items {
		// We can simply Fprint to a view.
		_, err := fmt.Fprintln(v, i.Name)
		if err != nil {
			log.Println("Error writing to the tasks view:", err)
		}
	}
	if len(items) != 0 {
		_, cy := v.Cursor()
		l, _ := v.Line(cy)
		models.CurrentTask = models.GetTask(l)
	}
	entriesView, _ := g.View(E)
	redrawEntries(g, entriesView)
	entriesView.Highlight = false
}

// Get the view and redraw it with current database info.
func redrawEntries(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	items := models.AllEntries(models.CurrentTask)
	if cv := g.CurrentView(); cv != nil && cv.Name() != P {
		for _, i := range items {
			// We can simply Fprint to a view.
			_, err := fmt.Fprintln(v, i.Name)
			// log.Println("redrawEntries CurrentTask.ID:", models.CurrentTask.ID)
			if err != nil {
				log.Println("Error writing to the entries view:", err)
			}
		}
	}
	// Set cursor to last item if we just finished inputting an entry.
	if bottom == false {
		v.SetCursor(0, 0)
	} else {
		v.SetCursor(0, len(items)-1)
	}
	// Read the cursor, get the data, print it.
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	if l != "" {
		models.CurrentEntry = models.GetEntry(l)
	}
	outputView, _ := g.View(O)
	redrawOutput(g, outputView)
}

// Start a new entry.
// This program does not use a timer, just timestamps.
func newEntry(g *gocui.Gui, v *gocui.View) error {
	var err error
	v.Highlight = false
	now := time.Now()
	e := models.StartEntry(models.CurrentTask, now)
	models.CurrentEntry = e
	ov, err := g.SetCurrentView("output")
	if err != nil {
		return err
	}
	ov.Clear()
	ov.Editable = true
	g.Cursor = true
	ov.SetCursor(0, 0)
	return err
}

// v will always equal output view
// This saves the text in the output view and does what it needs to do
// with it depending on what view it was called from (before switching to the output view).
func save(g *gocui.Gui, v *gocui.View) error {
	var err error
	if dView != nil && dView.Name() == P {
		// Project
		d := v.Buffer()
		models.CurrentProject.Description = d
		models.DB.Save(&models.CurrentProject)
		projectsView, _ := g.View(P)
		g.SetCurrentView(P)
		redrawProjects(g, projectsView)
		v.Editable = false
		g.Cursor = false
	} else if dView != nil && dView.Name() == T {
		// Task
		d := v.Buffer()
		models.CurrentTask.Description = d
		models.DB.Save(&models.CurrentTask)
		tasksView, _ := g.View(T)
		g.SetCurrentView(T)
		redrawTasks(g, tasksView)
		v.Editable = false
		g.Cursor = false
	} else {
		// Entry
		d := v.Buffer()
		models.StopEntry(models.CurrentEntry, time.Now(), d)
		entriesView, _ := g.View(E)
		entriesView.Highlight = true
		// bottom = true
		g.SetCurrentView(E)
		redrawEntries(g, entriesView)
		v.Editable = false
		g.Cursor = false
		bottom = false
	}
	dView = nil
	return err
}

// Get the view and redraw it with current database info.
// The output view should not need to be redrawn while it is itself selected,
// but we'll see...
// v is always the output view.
func redrawOutput(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	if cv := g.CurrentView(); cv != nil {
		// Projects
		if cv.Name() == P {
			h, m := models.CurrentProject.HoursMinutes()
			if _, err := fmt.Fprintf(v, "%d Hours\n%d Minutes\n\n",
				h, m); err != nil {
				log.Println("Error writing project time to the output view:", err)
			}
			if _, err := fmt.Fprintf(v, models.CurrentProject.Description); err != nil {
				log.Println("Error writing project description to the output view:", err)
			}
		}
		if cv.Name() == T {
			h, m := models.CurrentTask.HoursMinutes()
			if _, err := fmt.Fprintf(v, "%d Hours\n%d Minutes\n\n",
				h, m); err != nil {
				log.Println("Error writing task time to the output view:", err)
			}
			if _, err := fmt.Fprintf(v, models.CurrentTask.Description); err != nil {
				log.Println("Error writing task description to the output view:", err)
			}
		}
		// Entries
		if models.CurrentEntry.Start.IsZero() == false &&
			cv.Name() == E {
			details := models.CurrentEntry.Details
			start := models.CurrentEntry.Start.Format(models.TL)
			end := models.CurrentEntry.End.Format(models.TL)
			hours, minutes := models.CurrentEntry.HoursMinutes()

			if _, err := fmt.Fprintf(v, "%d Hours\n%d Minutes\nStart: %v\nEnd:   %v\n\n",
				hours, minutes, start, end); err != nil {
				log.Println("Error writing entry to the output view:", err)
			}
			if _, err := fmt.Fprintln(v, details); err != nil {
				log.Println("Error writing entry to the output view:", err)
			}
		}
	}
}

// This adds a description for projects and tasks.
func addDescription(g *gocui.Gui, v *gocui.View) error {
	var err error
	dView = v
	switch v.Name() {
	case P:
		ov, _ := g.SetCurrentView("output")
		ov.Clear()
		if _, err = fmt.Fprintf(ov, models.CurrentProject.Description); err != nil {
			log.Println("Error writing project description to the output view:", err)
		}
		ov.Editable = true
		g.Cursor = true
		ov.SetCursor(0, 0)
	case T:
		ov, _ := g.SetCurrentView("output")
		ov.Clear()
		if _, err = fmt.Fprintf(ov, models.CurrentTask.Description); err != nil {
			log.Println("Error writing task description to the output view:", err)
		}
		ov.Editable = true
		g.Cursor = true
		ov.SetCursor(0, 0)
	}
	return err
}

// The layout handler calculates all sizes depending on the current terminal size.
func layout(g *gocui.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()
	// Update the views according to the new terminal size.
	// Projects.
	_, err := g.SetView(P, 0, 0, pwidth, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update projects view")
	}
	// Tasks
	_, err = g.SetView(T, pwidth+1, 0, twidth, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update tasks view")
	}
	// Entries
	_, err = g.SetView(E, twidth+1, 0, ewidth, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update entries view")
	}
	// Output
	_, err = g.SetView("output", ewidth+1, 0, tw-1, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update output view")
	}
	// Status
	// Not used right now. If uncommented set all above SetView() y1 values to 'th-4'.
	// _, err = g.SetView("status", 0, th-sheight, tw-1, th-1)
	// if err != nil {
	// 	return errors.Wrap(err, "Cannot update input view.")
	// }
	return nil
}

// quit is a handler that gets bound to Ctrl-gocui. It signals the main loop to exit.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
