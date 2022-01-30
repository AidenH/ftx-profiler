package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

type ProgramSettings struct {
	VolumeSymbol string
}

var PrecisionMap = map[int]int{
	1: 10,
	2: 100,
	3: 1000,
	4: 10000,
	5: 100000,
	6: 1000000,
}

var ProfileUnitDiv = 10

func Round(input float64, precision int) float64 {

	var p int

	if precision == 0 {
		p = 1
	} else {
		s := []string{"1", strings.Repeat("0", precision)}
		p, _ = strconv.Atoi(strings.Join(s, ""))
	}

	pfloat := float64(p)

	result := math.Round(input*pfloat) / pfloat

	return result
}

// GuiDebugPrint prints debug strings to the selected Gui View
func GuiDebugPrint(v string, msg interface{}) error {

	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(v)
		if err != nil {
			return err
		}

		fmt.Fprintln(v, msg)

		return nil
	})

	return nil
}

func AddVData(price float64, size float64) error {

	VData[price] += Round(size, int(State.SizeGranularity))

	return nil
}
