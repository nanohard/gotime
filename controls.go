package main

import "github.com/jroimartin/gocui"

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", gocui.KeyCtrlA, gocui.ModNone, addItem); err != nil {
		return err
	}
	// if err := g.SetKeybinding("projects", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, copyInput); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }

	//
	// if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
	// 	return err
	// }

	return nil
}
