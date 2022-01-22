package main

import (
	"math"
	"strconv"
	"strings"
)

var PrecisionMap = map[int]int{
	1: 10,
	2: 100,
	3: 1000,
	4: 10000,
	5: 100000,
	6: 1000000,
}

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
