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
	PriceTrim      int
}

// InitCui initialize gocui cui
func InitCui() (*gocui.Gui, error) {

	// write cui defaults
	CState = CuiState{
		Middle:         0,
		SetMiddle:      true,
		LockWrite:      false,
		ProfileUnitDiv: 1,
		PriceTrim:      2,
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
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, Recenter); err != nil {
		panic(err)
	}

	// reset volume
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, ClearVolume); err != nil {
		panic(err)
	}

	// save volume data to file. utils.go: VolWrite
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, VolWrite); err != nil {
		panic(err)
	}

	initProfile("LOOKS-PERP", 0, 3, true, g)

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

	// set highlight color
	g.SelBgColor = gocui.ColorWhite
	g.SelFgColor = gocui.ColorBlue

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
		fMaxY := float64(maxY) // floatify Y-size for fewer casts during later for loop

		prec := float64(PrecisionMap[State.PricePrecision])
		// bring price up to int precision so we can subtract x amount of price levels
		// where x = half of the viewport rows. basically calculating the beginning of
		// the price ladder relative to global price precision
		modPrice := CState.Middle * prec // e.g. 0.0289 * 1000
		// shift price a relative amount of precision decimal points, then subtract
		// half of viewport rows, then shift back to proper price size
		ladderStart := (modPrice - float64(maxY/2)) / prec

		v.Clear()

		for i := fMaxY; i > 0.0; i-- {
			current := ladderStart + (i / prec)
			p := strconv.FormatFloat(current, 'f', State.PricePrecision, 64)
			f, _ := strconv.ParseFloat(p, 64)
			sizewidth := int(VData[f]) / CState.ProfileUnitDiv

			// rescale profile bar proportions if vol length above half view width
			if sizewidth >= maxX/3 {
				CState.ProfileUnitDiv *= 2
			}

			// print profile. if i = current price, mark on ladder
			if f == State.LastPrice {
				if VData[f] > 0 && State.VolumeCounts {
					fmt.Fprintf(v, "%s%s%s  -  %s %g\n",
						Color.Purple,
						p,
						Color.Default,
						strings.Repeat(Settings.VolumeSymbol, sizewidth),
						VData[f])

				} else {
					fmt.Fprintf(v, "%s%s%s  -  %s\n",
						Color.Purple,
						p[CState.PriceTrim:],
						Color.Default,
						strings.Repeat(Settings.VolumeSymbol, sizewidth))
				}

				// if price reaches extremes of visible profile, reset middle coord
				if i < 3 || i > (fMaxY-3) {
					CState.SetMiddle = true
				}

			} else {
				if VData[f] > 0 && State.VolumeCounts {
					fmt.Fprintf(v, "%s  -  %s %g\n",
						p,
						strings.Repeat(Settings.VolumeSymbol, sizewidth),
						VData[f])
				} else {
					fmt.Fprintf(v, "%s  -  %s\n",
						p[CState.PriceTrim:],
						strings.Repeat(Settings.VolumeSymbol, sizewidth))
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
				str := fmt.Sprintf("%s %s - %s %s",
					Color.Green,
					p[CState.PriceTrim:],
					size,
					Color.Default)
				fmt.Fprintln(v, str)

			} else if side == "sell" {
				str := fmt.Sprintf("%s %s - %s %s",
					Color.Red,
					p[CState.PriceTrim:],
					size,
					Color.Default)
				fmt.Fprintln(v, str)

			} else {
				err = errors.New("PrintTape - invalid side type")

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
