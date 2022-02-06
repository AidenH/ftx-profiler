package main

import (
	"net/http"
)

type AccountState struct {
	Result struct {
		PosSize   float64 `json:"totalPositionSize"`
		Balance   float64 `json:"totalAccountValue"`
		Positions []map[string]interface{}
	}
}

var Account = new(AccountState)
var client = &http.Client{}

// GetAccountInfo() retrieves user's FTX account info via http API request.
// Function will return an *http.Response for testing purposes to be used
// in conjunction with ParseHttpResp()
func (a *AccountState) GetAccountInfo() error {
	req, err := http.NewRequest("GET", "https://ftx.com/api/account", nil)
	if err != nil {
		return err
	}

	signature, ts, err := CreateSignature("/api/account")
	if err != nil {
		return err
	}

	req.Header.Add("FTX-KEY", Api)
	req.Header.Add("FTX-SIGN", signature)
	req.Header.Add("FTX-TS", ts)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	a.ParseHttpResp(resp)

	return nil
}

// GetOpenPositions() returns a user's open exchange positions
func (a *AccountState) GetOpenPositions() error {
	req, err := http.NewRequest("GET", "https://ftx.com/api/positions", nil)
	if err != nil {
		return err
	}

	signature, ts, err := CreateSignature("/api/positions")
	if err != nil {
		return err
	}

	req.Header.Add("FTX-KEY", Api)
	req.Header.Add("FTX-SIGN", signature)
	req.Header.Add("FTX-TS", ts)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := a.ParseHttpResp(resp); err != nil {
		return err
	}

	return nil
}
