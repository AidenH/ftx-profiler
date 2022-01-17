package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func InitCui() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	initProfile("NEAR-PERP", 1, 3, true, g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return g
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	var yOffset int

	if maxY > 40 {
		yOffset = 4
	} else {
		yOffset = 0
	}

	// STATUS BAR
	if v, err := g.SetView("status", 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "status bar")
	}

	// PROFILE
	if v, err := g.SetView("profile", 0, yOffset, (maxX/3*2)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "profile")
	}

	// TAPE
	if v, err := g.SetView("tape", maxX/3*2, yOffset, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = true
	}

	return nil
}

func PrintTape(g *gocui.Gui, side string, price float64, size string) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("tape")
		if err != nil {
			return err
		}

		if side == "buy" {
			fmt.Fprintln(v, fmt.Sprintf("%s %g - %s %s",
				"\033[32m", price, size, "\033[0m"))
		} else if side == "sell" {
			fmt.Fprintln(v, fmt.Sprintf("%s %g - %s %s",
				"\033[31m", price, size, "\033[0m"))
		} else {
			err = errors.New("PrintTape - invalid side type")
			return err
		}

		return nil
	})

	return nil
}

func setStatus(state ProfileState) error {
	state.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("status")
		if err != nil {
			return err
		}

		v.Clear()

		fmt.Fprintln(v, state.Market)

		return nil
	})

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
