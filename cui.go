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
	Middle         float64
	SetMiddle      bool
	LockWrite      bool
	ProfileUnitDiv int
}

// InitCui initialize gocui cui
func InitCui() (*gocui.Gui, error) {

	// write cui defaults
	CState = CuiState{
		Middle:         0,
		SetMiddle:      true,
		LockWrite:      false,
		ProfileUnitDiv: 1,
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	// quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// recenter profile
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, Recenter); err != nil {
		panic(err)
	}

	// clear volume map
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, ClearVolume); err != nil {
		panic(err)
	}

	// save volume data to file
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, VolWrite); err != nil {
		panic(err)
	}

	initProfile("SLP-PERP", 0, 4, true, g)

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

	// ORDERS AND POSITIONS
	if _, err := g.SetView("orders", 0, yOffset, (maxX/15)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	// PROFILE
	if v, err := g.SetView("profile", maxX/15, yOffset, (maxX/3*2)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = false
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

// Recenter recenters profile on screen
func Recenter(g *gocui.Gui, v *gocui.View) error {
	CState.SetMiddle = true

	GuiDebugPrint("status", "\nResetting profile...")

	return nil
}

// ClearVolume will empty VData and reset profile
func ClearVolume(g *gocui.Gui, v *gocui.View) error {
	VData = make(map[float64]float64)
	CState.ProfileUnitDiv = 1

	Recenter(g, v)
	GuiDebugPrint("status", "Clearing volume data...")

	if len(VData) != 0 {
		err := errors.New("VData not cleared")
		return err
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

		for i := fMaxY; i > 0.0; i-- {
			current := ladderStart + (i / prec)
			p := strconv.FormatFloat(current, 'f', State.PricePrecision, 64)
			f, _ := strconv.ParseFloat(p, 64)
			sizewidth := int(VData[f]) / CState.ProfileUnitDiv

			// rescale profile proportions if vol length above half view width
			if sizewidth >= maxX/3 {
				CState.ProfileUnitDiv *= 2
			}

			// print profile. if i = current price, mark on ladder
			if f == State.LastPrice {

				if VData[f] > 0 && State.VolumeCounts {
					//fmt.Fprintln(v, "\033[35m", p, "\033[0m- ",
					//	strings.Repeat(Settings.VolumeSymbol, sizewidth), VData[f])
					fmt.Fprintf(v, "%s%s%s  -  %s %g\n",
						Color.Purple,
						p,
						Color.Default,
						strings.Repeat(Settings.VolumeSymbol, sizewidth),
						VData[f])

				} else {
					//fmt.Fprintln(v, "\033[35m", p, "\033[0m- ",
					//	strings.Repeat(Settings.VolumeSymbol, sizewidth))
					fmt.Fprintf(v, "%s%s%s  -  %s\n",
						Color.Purple,
						p,
						Color.Default,
						strings.Repeat(Settings.VolumeSymbol, sizewidth))
				}

				if i < 3 || i > (fMaxY-3) {
					CState.SetMiddle = true
				}

			} else {
				if VData[f] > 0 && State.VolumeCounts {
					fmt.Fprintln(v, p, " - ", strings.Repeat(Settings.VolumeSymbol, sizewidth), VData[f])
				} else {
					fmt.Fprintln(v, p, " - ", strings.Repeat(Settings.VolumeSymbol, sizewidth))
				}
			}
		}

		return nil
	})

	return nil
}

func PrintOrders() error {
	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("orders")
		if err != nil {
			return err
		}

		fmt.Fprintln(v, OpenOrders.Result)

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
				str := fmt.Sprintf("%s %s - %s %s", Color.Green, p, size, Color.Default)
				fmt.Fprintln(v, str)

				//FileWrite(str)
			} else if side == "sell" {
				str := fmt.Sprintf("%s %s - %s %s", Color.Red, p, size, Color.Default)
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

		fmt.Fprintf(v, "%s - %s  OPEN: %s    BALANCE: %.2f",
			State.Market, p, o, Account.Result.Balance)

		return nil
	})

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
