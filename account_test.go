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
	dat := new(AccountState)

	resp, err = GetAccountInfo()
	if err != nil {
		log.Panicln(err)
	}

	dat.ParseHttpResp(resp)

	fmt.Println(dat.Result.Balance)
}
