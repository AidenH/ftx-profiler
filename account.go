package main

import "net/http"

var client = &http.Client{}

func GetAccountInfo() error {
	req, err := http.NewRequest("GET", "https://ftx.com/api/account", nil)
	if err != nil {
		return err
	}

	req.Header.Add("FTX-KEY", "")
	req.Header.Add("FTX-SIGN", "")
	req.Header.Add("FTX-TS", "")

	return nil
}
