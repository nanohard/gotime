package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/nanohard/gotime/models"
	"github.com/pkg/errors"
)

// var CurrentProject models.Project
// var CurrentTask models.Task
// var CurrentEntry models.Entry

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	// Check to make sure data exists in the next line,
	// otherwise disallow scroll down.
	if v != nil && lineBelow(g, v) == true {
		v.MoveCursor(0, 1, false)
		_, cy := v.Cursor()
		// var nv *gocui.View
		if v.Name() == "projects" {
			nv, _ := g.View("tasks")
			// n, _ := v.Word(cx, cy)
			n, _ := v.Line(cy)
			// log.Println("cursorDown Line:", v.Buffer())
			log.Println("cursorDown Line:", n)
			// n = strings.TrimSpace(n)
			// log.Println("cursorDown Line:", n)
			models.CurrentProject = models.GetProject(n)
			// log.Println("cursorDown CurrentProject:", models.CurrentProject)
			redrawTasks(g, nv)
		} else if v.Name() == "tasks" {
			//redrawEntries
		} else if v.Name() == "entries" {
			//	redrawOutput
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, -1, false)
		_, cy := v.Cursor()
		// var nv *gocui.View
		if v.Name() == "projects" {
			nv, _ := g.View("tasks")
			n, _ := v.Line(cy)
			// n, _ := v.Word(cx, cy)
			// log.Println("cursorUp Buffer:", v.Buffer())
			log.Println("cursorUp Line:", n)
			// n = strings.TrimSpace(n)
			// log.Println("cursorUp Line:", n)
			models.CurrentProject = models.GetProject(n)
			// log.Println("cursorUp CurrentProject:", models.CurrentProject)
			redrawTasks(g, nv)
		} else if v.Name() == "tasks" {
			//redrawEntries
		} else if v.Name() == "entries" {
			//	redrawOutput
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

// func getLine(g *gocui.Gui, v *gocui.View) error {
// 	var l string
// 	var err error
//
// 	_, cy := v.Cursor()
// 	if l, err = v.Line(cy); err != nil {
// 		l = ""
// 	}
//
// 	maxX, maxY := g.Size()
// 	// If there is data, then show msg.
// 	if l != "" {
// 		if v, err := g.SetView("msg", maxX/2-20, maxY/2, maxX/2+30, maxY/2+2); err != nil {
// 			if err != gocui.ErrUnknownView {
// 				return err
// 			}
// 			fmt.Fprintln(v, l)
// 			if _, err := g.SetCurrentView("msg"); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}
// 	return nil
// }

// Copy the input view (iv) and handle it.
// Used to add project or task.
func copyInput(g *gocui.Gui, iv *gocui.View) error {
	var e error
	// We want to read the view’s buffer from the beginning.
	iv.Rewind()
	// Get the output view via its name.
	var ov *gocui.View
	// If there is text input then add the item,
	// else go back to the input view.
	switch iv.Name() {
	case "addProject":
		ov, _ = g.View("projects")
		if iv.Buffer() != "" {
			models.AddProject(iv.Buffer())
			redrawProjects(g, ov)
		} else {
			addItem(g, ov)
			return nil
		}
	case "addTask":
		ov, _ = g.View("tasks")
		if iv.Buffer() != "" {
			models.AddTask(iv.Buffer(), models.CurrentProject)
			redrawTasks(g, ov)
		} else {
			addItem(g, ov)
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
	if err := g.DeleteView(iv.Name()); err != nil {
		log.Print(err)
	}
	// redrawView(g, ov)
	// Set the view back.
	if _, err := g.SetCurrentView(ov.Name()); err != nil {
		log.Print(err)
	}
	return e
}

// func delInput(g *gocui.Gui, v *gocui.View) error {
// 	if err := g.DeleteView("msg"); err != nil {
// 		return err
// 	}
// 	if _, err := g.SetCurrentView("projects"); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func delMsg(g *gocui.Gui, v *gocui.View) error {
// 	if err := g.DeleteView("msg"); err != nil {
// 		return err
// 	}
// 	if _, err := g.SetCurrentView("projects"); err != nil {
// 		return err
// 	}
// 	return nil
// }

// Add item to the current view (cv) using the text from the input view (iv).
func addItem(g *gocui.Gui, cv *gocui.View) error {
	maxX, maxY := g.Size()
	var title string
	var name string
	switch cv.Name() {
	case "projects":
		title = "Name of new project"
		name = "addProject"
		// models.AddProject()
	case "tasks":
		title = "Name of new task"
		name = "addTask"
		// case "entries":
		//     title = "Name of new entry"
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

// // Get the current view (cv) and transfer cursor to the new view (nv).
// // If the new view has no items then prompt user to add one.
func selectItem(g *gocui.Gui, cv *gocui.View) error {
	// var e error
	// var nv string
	// _, cy := cv.Cursor()
	switch cv.Name() {
	case "projects":
		if _, err := g.SetCurrentView("tasks"); err != nil {
			return err
		}
		// models.CurrentProject, _ = cv.Line(cy)
	case "tasks":
		// nv = "entries"
		// case "entries":
		//     nv = "ad-hoc"
	}
	// v, err := g.SetCurrentView(nv)
	// if err != nil {
	// 	return err
	// }
	// redrawView(g, v)
	return nil
}

// Get the current view (cv) and transfer cursor to the new view (nv).
func goBack(g *gocui.Gui, cv *gocui.View) error {
	// var e error
	// var nv string
	switch cv.Name() {
	case "tasks":
		if _, err := g.SetCurrentView("projects"); err != nil {
			return err
		}
	case "entries":
		// nv = "tasks"
		// case "entries":
		//     nv = "ad-hoc"
	}
	// v, err := g.SetCurrentView(nv)
	// if err != nil {
	// 	return err
	// }
	// redrawView(g, v)
	return nil
}

// Get the projects view and redraw it with current database info.
func redrawProjects(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	items := models.AllProjects()
	// Loop through projects to add their names to the view.
	for _, i := range items {
		// Again, we can simply Fprint to a view.
		_, err := fmt.Fprintln(v, i.Name)
		if err != nil {
			log.Println("Error writing to the projects view:", err)
		}
	}
	if len(items) == 0 {
		addItem(g, v)
	}
	models.CurrentProject = models.GetProject(l)
}

// Get the view and redraw it with current database info.
func redrawTasks(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	items := models.AllTasks(models.CurrentProject)
	// Loop through tasks to add their names to the view.
	for _, i := range items {
		// Again, we can simply Fprint to a view.
		_, err := fmt.Fprintln(v, i.Name)
		if err != nil {
			log.Println("Error writing to the projects view:", err)
		}
	}
	if len(items) != 0 {
		models.CurrentTask = models.GetTask(l)
	}

}

// func redrawProjects(v *gocui.View) {
// 	// Projects
// 	projectItems := models.AllProjects()
// 	// panic(projectItems)
// 	// Loop through projects to add their names to the view.
// 	for _, p := range projectItems {
// 		// Again, we can simply Fprint to a view.
// 		_, err := fmt.Fprintln(v, p.Name)
// 		if err != nil {
// 			log.Println("Error writing to the projects view:", err)
// 		}
// 	}
// }

// The layout handler calculates all sizes depending on the current terminal size.
func layout(g *gocui.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()
	// Update the views according to the new terminal size.
	// Projects.
	_, err := g.SetView("projects", 0, 0, pwidth, th-4)
	if err != nil {
		return errors.Wrap(err, "Cannot update projects view")
	}
	// Tasks
	_, err = g.SetView("tasks", pwidth+1, 0, twidth, th-4)
	if err != nil {
		return errors.Wrap(err, "Cannot update tasks view")
	}
	// Entries
	_, err = g.SetView("entries", twidth+1, 0, ewidth, th-4)
	if err != nil {
		return errors.Wrap(err, "Cannot update entries view")
	}
	// Output
	_, err = g.SetView("output", ewidth+1, 0, tw-1, th-4)
	if err != nil {
		return errors.Wrap(err, "Cannot update output view")
	}
	// Status
	_, err = g.SetView("status", 0, th-sheight, tw-1, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update input view.")
	}
	return nil
}

// quit is a handler that gets bound to Ctrl-gocui. It signals the main loop to exit.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
