package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/nanohard/gotime/models"
	"github.com/pkg/errors"
)

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	// Check to make sure data exists in the next line,
	// otherwise disallow scroll down.
	if v != nil && lineBelow(g, v) == true {
		v.MoveCursor(0, 1, false)
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, -1, false)
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
	// We want to read the view’s buffer from the beginning.
	iv.Rewind()

	// Get the output view via its name.
	var ov *gocui.View
	var p models.Project
	switch iv.Name() {
	case "addProject":
		ov, _ = g.View("projects")
		p = models.AddProject(iv.Buffer())
	case "addTask":
		ov, _ = g.View("tasks")
	}
	_, e := fmt.Fprint(ov, p.Name)
	if e != nil {
		log.Panic("Cannot print to output view:", e)
	}
	// Thanks to views being an io.Writer, we can simply Fprint to a view.
	// _, e := fmt.Fprint(ov, iv.Buffer())
	// if e != nil {
	// 	log.Panic("Cannot print to output view:", e)
	// }
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
	if _, err := g.SetCurrentView(ov.Name()); err != nil {
		log.Print(err)
	}
	return e
	// return nil
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
