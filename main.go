package main

import (
	"github.com/jroimartin/gocui"
)

type ProfileState struct {
	Market          string
	SizeGranularity float64
	PricePrecision  int
	Aggregate       bool
	Gui             *gocui.Gui
}

var VData = make(map[float64]float64)

var State = ProfileState{}

// initProfile initializes an FTX websocket to populate tape and profile
// views
func initProfile(
	mar string,
	gran float64,
	price int,
	agg bool,
	g *gocui.Gui) error {

	State = ProfileState{
		Market:          mar,
		SizeGranularity: gran,
		PricePrecision:  price,
		Aggregate:       agg,
		Gui:             g,
	}

	err := SocketInit()

	setStatus()

	return err
}

func main() {

	_ = InitCui()

	//socket := initProfile("NEAR-PERP", 1, 0, true, g)
}
