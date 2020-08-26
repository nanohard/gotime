package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

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
	ewidth = 65
	// Input box height.
	sheight = 3
	// P is string for projects.
	P = "projects"
	// T is string for tasks.
	T = "tasks"
	// E is string for entries.
	E = "entries"
	// O is string for output.
	O = "output"
)

func main() {
	// Debug log
	usr, _ := user.Current()
	dir := usr.HomeDir
	f, err := os.OpenFile(dir+"/.gotime.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file no exists")
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
	g.SelFgColor = gocui.ColorBlue
	g.BgColor = gocui.ColorBlack
	g.FgColor = gocui.ColorWhite

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
	projectView.SelBgColor = gocui.ColorBlue
	projectView.SelFgColor = gocui.ColorBlack

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
	tasksView.SelBgColor = gocui.ColorBlue
	tasksView.SelFgColor = gocui.ColorBlack

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
	entriesView.SelBgColor = gocui.ColorBlue
	entriesView.SelFgColor = gocui.ColorBlack

	// Output view.
	outputView, err := g.SetView("output", ewidth+1, 0, terminalWidth-1, terminalHeight-4)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create output view (AAAGGHHH!!!):", err)
		return
	}
	outputView.FgColor = gocui.ColorWhite
	// Let the view scroll if the output exceeds the visible area.
	outputView.Autoscroll = true
	outputView.Wrap = true

	// Status view.
	// Not used right now.
	// statusView, err := g.SetView("status", 0, terminalHeight-sheight, terminalWidth-1, terminalHeight-1)
	// if err != nil && err != gocui.ErrUnknownView {
	// 	log.Println("Failed to create input view:", err)
	// 	return
	// }
	// statusView.Title = "Status"
	// statusView.FgColor = gocui.ColorYellow

	// Database ***************************************************
	models.InitDB()
	defer models.DB.Close()

	// Projects
	projectItems := models.AllProjects()
	if len(projectItems) > 0 {
		models.CurrentProject = projectItems[0]
		redrawProjects(g, projectView)
	}

	// Main loop stuff *********************************************
	// Apply keybindings to program.
	if err = keybindings(g); err != nil {
		log.Panicln(err)
	}
	// Must set initial view here, right before program start!!!
	v, _ := g.SetCurrentView(P)
	// Move the cursor to update the output view with the description.
	// (workaround)
	cursorUp(g, v)

	// If no projects on start then prompt the user to add a project.
	if len(projectItems) == 0 {
		inputView(g, projectView)
	}
	// Start the main loop.
	err = g.MainLoop()
	log.Println("Main loop has finished:", err)
}
