package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	VolMinFilter float64
}

var VData = make(map[float64]float64)
var Ladder = make(map[float64]int)

var OpenOrders = OrdersRestReply{}
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
		VolMinFilter:    10,
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
	ts := time.Now().Format(time.Stamp)
	var err error

	// check application dirs
	MakeDirs()

	// create session log file
	LogFile, err = os.Create(fmt.Sprintf("%s/.ftx-profiler/profiler-output-log-%s",
		HomeDir,
		strings.Replace(ts, " ", "-", -1)),
	)
	if err != nil {
		log.Println("unable to create log file")
		panic(err)
	}
	defer LogFile.Close()

	_, err = InitCui()
	if err != nil {
		FileWrite("cui main loop error")
		FileWrite(err.Error())
		panic(err)
	}

	LogFile.Sync()

	// ideally initProfile would be started here with a passed gui instance, instead
	// of at cui.go:61. but still have not been able to make this work
	//socket := initProfile("NEAR-PERP", 1, 0, true, g)
}
