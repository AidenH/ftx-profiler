package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/sacOO7/gowebsocket"
)

type ProgramState struct {
	Connections []Connection

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

type Connection struct {
	Name      string
	Socket    gowebsocket.Socket
	Subscribe func(Connection) error
}

var VData = make(map[float64]float64)
var Ladder = make(map[float64]int)

var OpenOrders = OrdersRestReply{}
var State = ProgramState{}
var Account = AccountState{
	Orders: make(map[int]Order), // init orders map
}
var CState = CuiState{}

var client = &http.Client{}

var HomeDir, _ = os.UserHomeDir()
var LogFile *os.File
var VolFile *os.File

// initProfile initializes an FTX websocket to populate tape and profile
// views
func initProfile(g *gocui.Gui) error {

	State = ProgramState{
		Market:          Config.Market,
		SizeGranularity: Config.SizeGranularity,
		PricePrecision:  Config.PricePrecision,
		Aggregate:       Config.Aggregate,
		Gui:             g,
		TapeTrue:        true,
		ProfileTrue:     true,
		VolMinFilter:    10,
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
