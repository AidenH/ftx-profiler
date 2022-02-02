package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"math"
	"os"
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

func Round(input float64, precision int) (float64, error) {
	//FileWrite(fmt.Sprintf("Round\ninput:%f", input))

	var p int
	var err error

	if precision == 0 {
		p = 1
	} else {
		s := []string{"1", strings.Repeat("0", precision)}
		p, err = strconv.Atoi(strings.Join(s, ""))
		if err != nil {
			return 0.0, err
		}
	}

	pfloat := float64(p)

	result := math.Round(input*pfloat) / pfloat

	//FileWrite(fmt.Sprintf("Rounded: %f\n", result))

	return result, nil
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

func FileWrite(dat string) error {
	_, err := LogFile.Write([]byte(fmt.Sprint(dat, "\n")))
	if err != nil {
		return err
	}

	return nil
}

func AddVData(price float64, size float64) error {

	out, err := Round(size, int(State.SizeGranularity))
	if err != nil {
		return err
	}

	VData[price] += out

	return nil
}

func CreateSignature(msg string) (hash.Hash, error) {
	//ts := time.Now().Unix()
	secret := "test"
	data := "test"

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))

	return nil, nil
}
