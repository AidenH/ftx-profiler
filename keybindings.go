package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func Keybindings(g *gocui.Gui) {
	// quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// recenter profile
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, Recenter); err != nil {
		panic(err)
	}

	// reset volume
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, ClearVolume); err != nil {
		panic(err)
	}

	// save volume data to file. utils.go: VolWrite()
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, VolWrite); err != nil {
		panic(err)
	}

}
