package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jroimartin/gocui"
)

type CuiState struct {
	Middle    float64
	SetMiddle bool
}

// InitCui initialize gocui cui
func InitCui() *gocui.Gui {

	CState = CuiState{
		Middle:    0,
		SetMiddle: true,
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	initProfile("NEAR-PERP", 0, 2, true, g)

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

		v.Wrap = true

		fmt.Fprintln(v, "init profile")

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

// PrintProfile writes volume profile data to its view
func PrintProfile() error {
	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("profile")
		if err != nil {
			return err
		}

		_, maxY := v.Size()
		fMaxY := float64(maxY)

		prec := float64(PrecisionMap[State.PricePrecision])
		modPrice := State.LastPrice * prec
		ladderStart := (modPrice - float64(maxY/2)) / prec

		v.Clear()

		// DEBUG output VData
		//fmt.Fprintln(v, VData)

		//for i := 0.0; i < fMaxY; i++ {
		for i := fMaxY; i > 0.0; i-- {
			current := ladderStart + (i / prec)

			p := strconv.FormatFloat(current, 'f', State.PricePrecision, 64)

			if current == State.LastPrice {
				fmt.Fprintln(v, "\033[35m", p, "\033[0m - ", VData[current])
			} else {
				fmt.Fprintln(v, p, " - ", VData[current])
			}

		}

		return nil
	})

	return nil
}

// PrintTape outputs formulated trade event strings to the tape view
func PrintTape(side string, price float64, size string) error {
	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("tape")
		if err != nil {
			return err
		}

		dt := time.Now()
		fmt.Fprint(v, dt.Format(time.StampMicro)[7:19])

		p := strconv.FormatFloat(price, 'f', State.PricePrecision, 64)

		if side == "buy" {
			fmt.Fprintln(v, fmt.Sprintf("%s %s - %s %s",
				"\033[32m", p, size, "\033[0m"))
		} else if side == "sell" {
			fmt.Fprintln(v, fmt.Sprintf("%s %s - %s %s",
				"\033[31m", p, size, "\033[0m"))
		} else {
			err = errors.New("PrintTape - invalid side type")
			return err
		}

		return nil
	})

	return nil
}

// SetStatus updates the status bar with relevant market information
func SetStatus() error {
	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("status")
		if err != nil {
			return err
		}

		v.Clear()

		p := strconv.FormatFloat(State.LastPrice, 'f', State.PricePrecision, 64)

		fmt.Fprintf(v, "%s - %s", State.Market, p)

		return nil
	})

	return nil
}

// GuiDebugPrint prints debug strings to the selected Gui View
func GuiDebugPrint(v string, msg string) error {
	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(v)
		if err != nil {
			return err
		}

		fmt.Fprint(v, msg)

		return nil
	})

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
