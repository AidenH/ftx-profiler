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

func TestRetrieveAccountInfo(t *testing.T) {
	go RetrieveAccountInfo()
}
