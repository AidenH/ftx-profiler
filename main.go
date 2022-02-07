package main

import (
	"net/http"
	"os"

	"github.com/jroimartin/gocui"
)

type ProgramState struct {
	Market          string
	LastPrice       float64
	OpenPrice       float64
	SizeGranularity float64
	PricePrecision  int
	Aggregate       bool
	Gui             *gocui.Gui
	TapeTrue        bool
	ProfileTrue     bool
}

var VData = make(map[float64]float64)
var Ladder = make(map[float64]int)

var OpenOrders = Orders{}
var Settings = ProgramSettings{}
var State = ProgramState{}
var Account = AccountState{}
var CState = CuiState{}

var client = &http.Client{}

var LogFile, _ = os.Create("/home/lurkcs/profiler-output-log")

// initProfile initializes an FTX websocket to populate tape and profile
// views
func initProfile(
	mar string,
	gran float64,
	price int,
	agg bool,
	g *gocui.Gui) error {

	State = ProgramState{
		Market:          mar,
		SizeGranularity: gran,
		PricePrecision:  price,
		Aggregate:       agg,
		Gui:             g,
		TapeTrue:        true,
		ProfileTrue:     true,
	}

	Settings = ProgramSettings{
		// Recommend '#' or '█'
		VolumeSymbol: "█",
	}

	HandleOsArgs()

	// init continuous account info checks
	go RetrieveAccountInfo()

	err := SocketInit()

	return err
}

func main() {
	defer LogFile.Close()

	InitCui()

	LogFile.Sync()

	//socket := initProfile("NEAR-PERP", 1, 0, true, g)
}
