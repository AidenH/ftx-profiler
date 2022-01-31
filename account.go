package main

import (
	"net/http"
)

var client = &http.Client{}

func GetAccountInfo() (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://ftx.com/api/account", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("FTX-KEY", "")
	req.Header.Add("FTX-SIGN", "")
	req.Header.Add("FTX-TS", "")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
