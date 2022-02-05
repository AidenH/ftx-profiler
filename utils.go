package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type ProgramSettings struct {
	VolumeSymbol string
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

func Recenter(*gocui.Gui, *gocui.View) error {
	GuiDebugPrint("status", "\nResetting profile...")
	CState.SetMiddle = true

	return nil
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

func AddVData(price float64, size float64) error {

	/*out, err := Round(size, int(State.SizeGranularity))
	if err != nil {
		return err
	}

	GuiDebugPrint("tape", fmt.Sprint(price, " ", size))

	VData[float64(price)] += out*/

	VData[price] += size

	return nil
}

func CreateSignature(msg string) (string, string, error) {
	ts := time.Now().UnixMilli()
	data := fmt.Sprint(ts, "GET", msg)

	h := hmac.New(sha256.New, []byte(ApiSecret))
	h.Write([]byte(data))

	sha := hex.EncodeToString(h.Sum(nil))

	return sha, fmt.Sprint(ts), nil
}

func (a *AccountState) ParseHttpResp(resp *http.Response) error {
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rbody, &a); err != nil {
		return err
	}

	return nil
}
