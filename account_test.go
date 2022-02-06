package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

var resp *http.Response
var err error = nil

func TestGetAccountInfo(t *testing.T) {
	var Account = new(AccountState)

	if err := Account.GetAccountInfo(); err != nil {
		log.Panicln(err)
	}

	fmt.Print("\nAccount Info: ", Account.Result, "\n\n")
}

// not sure this is necessary as GetAccountInfo gets all the relevant info
// in a more consistent manner
func TestGetOpenPositions(t *testing.T) {
	Account := new(AccountState)

	if err := Account.GetOpenPositions(); err != nil {
		log.Panicln(err)
	}

	fmt.Println("Positions: ", Account)
}
