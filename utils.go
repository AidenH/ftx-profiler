package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type ProgramSettings struct {
	VolumeSymbol string
	PriceMarker  string
}

var Color = struct {
	Green   string
	Red     string
	Purple  string
	Default string
}{
	Green:   "\033[32m",
	Red:     "\033[31m",
	Purple:  "\033[35m",
	Default: "\033[0m",
}

// Negative price precisions not working currently - so no BTC as of yet!
var PrecisionMap = map[int]int{
	-3: -1000,
	-2: -100,
	-1: -10,
	0:  1,
	1:  10,
	2:  100,
	3:  1000,
	4:  10000,
	5:  100000,
	6:  1000000,
}

func Round(input float64, precision int) (float64, error) {
	//FileWrite(fmt.Sprintf("Round\ninput:%f", input))

	var p int
	var err error

	if precision == 0 {
		p = 1
	} else if precision > 0 {
		s := []string{"1", strings.Repeat("0", precision)}
		p, err = strconv.Atoi(strings.Join(s, ""))
		if err != nil {
			return 0.0, err
		}
	} else {
		p = PrecisionMap[precision]
	}

	pfloat := float64(p)

	result := math.Round(input*pfloat) / pfloat

	return result, nil
}

// GuiDebugPrint prints debug strings to the selected Gui View
func GuiDebugPrint(view string, msg interface{}) error {

	State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}

		fmt.Fprintln(v, msg)

		return nil
	})

	return nil
}

func FileWrite(dat string) error {
	_, err := LogFile.Write([]byte(fmt.Sprint(dat, "\n")))
	if err != nil {
		return err
	}

	return nil
}

// VolWrite saves volume data from VData map to ~/.ftx-profiler/voldata-<date>
func VolWrite(*gocui.Gui, *gocui.View) error {
	ts := time.Now().Format(time.Stamp)

	// declare file name
	filename := fmt.Sprintf("%s/.ftx-profiler/%s-voldata-%s",
		HomeDir,
		State.Market,
		strings.Replace(ts, " ", "-", -1))
	var err error

	// file var
	VolFile, err = os.Create(filename)
	if err != nil {
		panic(err)
	}

	e := gob.NewEncoder(VolFile)

	if err = e.Encode(VData); err != nil {
		panic(err)
	}

	VolFile.Close()

	return nil
}

func VolRead(f string) error {
	filename, err := os.Open(f)
	if err != nil {
		return err
	}

	defer filename.Close()

	d := gob.NewDecoder(filename)

	if err := d.Decode(&VData); err != nil {
		return err
	}

	FileWrite(fmt.Sprint("VData:\n", VData))

	return nil
}

func AddVData(price float64, size float64) error {

	VData[price] += size

	return nil
}

func CreateHttpSignature(msg string) (string, string, error) {
	ts := time.Now().UnixMilli()
	data := fmt.Sprint(ts, "GET", msg)

	h := hmac.New(sha256.New, []byte(ApiSecret))
	h.Write([]byte(data))

	sha := hex.EncodeToString(h.Sum(nil))

	return sha, fmt.Sprint(ts), nil
}

func CreateSocketSignature() (string, int64, error) {
	ts := time.Now().UnixMilli()
	data := fmt.Sprint(ts, "websocket_login")

	h := hmac.New(sha256.New, []byte(ApiSecret))

	_, err := h.Write([]byte(data))
	if err != nil {
		return "", 0, err
	}

	sha := hex.EncodeToString(h.Sum(nil))

	return sha, ts, nil
}

func ParseHttpResp(resp *http.Response, out interface{}) error {
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in ParseHttpResp: ioutil:")
		fmt.Println("resp.Body output: ", resp.Body)
		return err
	}

	// apply json data to relevant Account struct fields
	if err := json.Unmarshal(rbody, &out); err != nil {
		fmt.Println("error in ParseHttpResp: unmarshal:")
		fmt.Println("rbody output: ", string(rbody))
		return err
	}

	return nil
}

// RetrieveAccountInfo, called as a goroutine, will perform intermittent checks on user's
// account balance, open positions and other account information. To be added to as
// necessary
func RetrieveAccountInfo() {
	for {

		// http request get account info
		if err := Account.GetAccountInfo(); err != nil {
			log.Panicln(err)
		}

		o := Account.Position

		// fill out active account information
		for _, i := range Account.Result.PositionsData {
			if i.Future == State.Market {
				if i.Size > 0 {
					o.Entry = i.Entry
					o.Size = i.Size
					o.Side = i.Side
					o.Pnl = i.Pnl
				} else {
					o.Entry = 0
					o.Size = 0
					o.Side = ""
					o.Pnl = 0
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

// MakeDirs checks for and makes application directories if necessary
func MakeDirs() {
	if err := os.Mkdir(fmt.Sprintf("%s/.ftx-profiler", HomeDir), 0700); err != nil {
		errType := fmt.Sprintf("%T", err)

		// does ~/.ftx-profiler/ already exist?
		if errType == "*fs.PathError" {
			log.Println("~/.ftx-profiler/ already present")
		} else {
			panic(err)
		}
	}

	if err := os.Mkdir(fmt.Sprintf("%s/.ftx-profiler/vol", HomeDir), 0700); err != nil {
		errType := fmt.Sprintf("%T", err)

		// does ~/.ftx-profiler/ already exist?
		if errType == "*fs.PathError" {
			log.Println("~/.ftx-profiler/ already present")
		} else {
			panic(err)
		}
	}

}
