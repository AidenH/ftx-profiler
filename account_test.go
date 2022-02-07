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
	if err := Account.GetAccountInfo(); err != nil {
		log.Println(err)
		t.Fail()
	}

	fmt.Print("\nAccount Info: ", Account.Result, "\n\n")
}

func TestGetOpenOrders(t *testing.T) {
	if err := OpenOrders.GetOpenOrders(); err != nil {
		log.Println(err)
		t.Fail()
	}

	fmt.Print("\nOpen Orders: ", OpenOrders, "\n\n")
}

// unifinished test
func TestRetrieveAccountInfo(t *testing.T) {
	go RetrieveAccountInfo()
}
