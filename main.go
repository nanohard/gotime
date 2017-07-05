package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/nanohard/gotime/models"
)

const (
	// The following 3 boxes will allow for 21 viewable characters.
	// They will not adjust horizontally, only vertically,
	// so only the output box will readjust both horizontally and vertically.
	//
	// Projects box width.
	pwidth = 22
	// Tasks box width.
	twidth = 44
	// Entries box width.
	ewidth = 61
	// Input box height.
	sheight = 3
	// P is string for projects.
	P = "projects"
	// T is string for tasks.
	T = "tasks"
	// E is string for entries.
	E = "entries"
)

func main() {
	// Debug log
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Println("Failed to create a GUI:", err)
		return
	}
	defer g.Close()

	// Highlight active view.
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	// The GUI object wants to know how to manage the layout.
	// Unlike termui, gocui does not use a grid layout.
	// Instead, it relies on a custom layout handler function to manage the layout.
	//
	// Here we set the layout manager to a function named layout that is defined further down.
	g.SetManagerFunc(layout)

	// Bind the quit handler function (also defined further down) to Ctrl-C,
	// so that we can leave the application at any time.
	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	// View definitions *******************************************************************
	// The terminal’s width and height are needed for layout calculations.
	terminalWidth, terminalHeight := g.Size()
	// Projects view.
	projectView, err := g.SetView(P, 0, 0, pwidth, terminalHeight-4)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create projects view:", err)
		return
	}
	projectView.Title = "Projects"
	projectView.FgColor = gocui.ColorCyan
	projectView.Highlight = true
	projectView.SelBgColor = gocui.ColorGreen
	projectView.SelFgColor = gocui.ColorBlack
	// projectView.Editable = true

	// Tasks view.
	tasksView, err := g.SetView(T, pwidth+1, 0, twidth, terminalHeight-4)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return
	}
	tasksView.Title = "Tasks"
	tasksView.FgColor = gocui.ColorCyan

	// // Entries view.
	entriesView, err := g.SetView(E, twidth+1, 0, ewidth, terminalHeight-4)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create main view:", err)
		return
	}
	entriesView.Title = "Entries"
	entriesView.FgColor = gocui.ColorCyan

	// Output view.
	outputView, err := g.SetView("output", ewidth+1, 0, terminalWidth-1, terminalHeight-4)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create output view (AAAGGHHH!!!):", err)
		return
	}
	outputView.FgColor = gocui.ColorGreen
	// Let the view scroll if the output exceeds the visible area.
	outputView.Autoscroll = true
	_, err = fmt.Println(outputView, "Press Ctrl-c to quit")
	if err != nil {
		log.Println("Failed to print into output view (2):", err)
	}
	outputView.Wrap = true

	// Status view.
	statusView, err := g.SetView("status", 0, terminalHeight-sheight, terminalWidth-1, terminalHeight-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create input view:", err)
		return
	}
	statusView.Title = "Status"
	statusView.FgColor = gocui.ColorYellow

	// The input view shall be editable.
	// inputView.Editable = true
	// err = inputView.SetCursor(0, 0)
	// if err != nil {
	// 	log.Println("Failed to set cursor:", err)
	// 	return
	// }
	// Set the focus to the input view.
	// _, err = g.SetCurrentView("input")
	// Activate the cursor for the current view.
	// g.Cursor = true
	// if err != nil {
	// 	log.Println("Cannot set focus to input view:", err)
	// }

	// Database ***************************************************
	models.InitDB()
	defer models.DB.Close()

	// Projects
	projectItems := models.AllProjects()
	if len(projectItems) > 0 {
		models.CurrentProject = projectItems[0]
		redrawProjects(g, projectView)
		redrawTasks(g, tasksView)
	}

	// panic(projectItems)
	// Loop through projects to add their names to the view.
	// for _, p := range projectItems {
	// 	// Again, we can simply Fprint to a view.
	// 	_, err = fmt.Fprintln(projectView, p.Name)
	// 	if err != nil {
	// 		log.Println("Error writing to the projects view:", err)
	// 	}
	// }
	// If there is at least one project item then it will be highlighted on program start.
	// So then we also want to load its tasks to be shown on start as well.
	// if len(projectItems) > 0 {
	// 	//models.CurrentProject = projectItems[0]
	// 	redrawTasks(g, tasksView)
	// }

	// Main loop stuff *********************************************
	// Apply keybindings to program.
	if err = keybindings(g); err != nil {
		log.Panicln(err)
	}
	// Must set initial view here, right before program start!!!
	if _, err = g.SetCurrentView(P); err != nil {
		log.Panic(err)
	}
	// If no projects on start then prompt the user to add a project.
	if len(projectItems) == 0 {
		inputView(g, projectView)
	}
	// Start the main loop.
	err = g.MainLoop()
	log.Println("Main loop has finished:", err)
}
