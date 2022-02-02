package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type CuiState struct {
	Middle    float64
	SetMiddle bool
}

// InitCui initialize gocui cui
func InitCui() (*gocui.Gui, error) {

	CState = CuiState{
		Middle:    0,
		SetMiddle: true,
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	initProfile("BTC-PERP", 0, 0, true, g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return nil, err
	}

	return g, err
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
	if _, err := g.SetView("status", 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	// PROFILE
	if v, err := g.SetView("profile", 0, yOffset, (maxX/3*2)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
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

// PrintProfile writes profile with volume data to the profile view
func PrintProfile() error {
	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("profile")
		if err != nil {
			FileWrite("err getting PrintProfile view")
			return err
		}

		maxX, maxY := v.Size()
		fMaxY := float64(maxY)

		prec := float64(PrecisionMap[State.PricePrecision])
		modPrice := CState.Middle * prec
		ladderStart := (modPrice - float64(maxY/2)) / prec

		v.Clear()
		//Ladder = make(map[float64]int)

		for i := fMaxY; i > 0.0; i-- {
			current := ladderStart + (i / prec)
			p := strconv.FormatFloat(current, 'f', State.PricePrecision, 64)
			f, _ := strconv.ParseFloat(p, 64)
			sizewidth := int(VData[f]) / ProfileUnitDiv

			// rescale profile proportions if vol length above half view width
			if sizewidth >= maxX/3 {
				ProfileUnitDiv *= 2
			}

			// print profile. if i = current price, mark on ladder
			if f == State.LastPrice {
				fmt.Fprintln(v, "\033[35m", p, "\033[0m- ", strings.Repeat(Settings.VolumeSymbol, sizewidth))

				if i < 3 || i > (fMaxY-3) {
					CState.SetMiddle = true
				}

			} else {
				fmt.Fprintln(v, p, " - ", strings.Repeat(Settings.VolumeSymbol, sizewidth))
			}
		}

		return nil
	})

	return nil
}

// PrintTape outputs formulated trade event strings to the tape view
func PrintTape(side string, price float64, size string) error {
	if State.TapeTrue {
		State.Gui.Update(func(g *gocui.Gui) error {
			v, err := g.View("tape")
			if err != nil {
				return err
			}

			dt := time.Now()
			fmt.Fprint(v, dt.Format(time.StampMicro)[7:19])

			p := strconv.FormatFloat(price, 'f', State.PricePrecision, 64)

			if side == "buy" {
				str := fmt.Sprintf("%s %s - %s %s", "\033[32m", p, size, "\033[0m")
				fmt.Fprintln(v, str)

				//FileWrite(str)
			} else if side == "sell" {
				str := fmt.Sprintf("%s %s - %s %s", "\033[31m", p, size, "\033[0m")
				fmt.Fprintln(v, str)

				//FileWrite(str)
			} else {
				err = errors.New("PrintTape - invalid side type")

				//FileWrite(side)
				return err
			}

			return nil
		})
	}

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
		o := strconv.FormatFloat(State.OpenPrice, 'f', State.PricePrecision, 64)

		fmt.Fprintf(v, "%s - %s  OPEN: %s", State.Market, p, o)

		return nil
	})

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
