package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/gocui"
)

type ProgramState struct {
	// market
	Market          string
	LastPrice       float64
	OpenPrice       float64
	SizeGranularity float64
	PricePrecision  int
	Aggregate       bool

	// application
	Gui          *gocui.Gui
	TapeTrue     bool
	ProfileTrue  bool
	VolumeCounts bool
}

var VData = make(map[float64]float64)
var Ladder = make(map[float64]int)

var OpenOrders = Orders{}
var Settings = ProgramSettings{}
var State = ProgramState{}
var Account = AccountState{}
var CState = CuiState{}

var client = &http.Client{}

var HomeDir, _ = os.UserHomeDir()
var LogFile *os.File
var VolFile *os.File

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
		// i recommend '#' or '█'
		VolumeSymbol: "█",
		PriceMarker:  "<",
	}

	HandleOsArgs()

	// init continuous account info checks
	go RetrieveAccountInfo()

	err := SocketInit()

	return err
}

func main() {
	var err error

	if err = os.Mkdir(fmt.Sprintf("%s/.ftx-profiler", HomeDir), 0700); err != nil {
		errType := fmt.Sprintf("%T", err)

		// does ~/.ftx-profiler/ already exist?
		if errType == "*fs.PathError" {
			log.Println("~/.ftx-profiler/ already present")
		} else {
			panic(err)
		}
	}

	LogFile, err = os.Create(fmt.Sprintf("%s/.dwm/profiler-output-log", HomeDir))
	if err != nil {
		log.Println("unable to create log file")
		panic(err)
	}
	defer LogFile.Close()

	InitCui()

	LogFile.Sync()

	// ideally initProfile would be started here with a passed gui instance, instead
	// of at cui.go:61. but still have not been able to make this work
	//socket := initProfile("NEAR-PERP", 1, 0, true, g)
}
