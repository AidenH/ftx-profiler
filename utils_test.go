package main

import (
	"fmt"
	"testing"
)

func TestRound(t *testing.T) {
	input := 19.5955

	result, _ := Round(input, 0)
	fmt.Println(result)
}
