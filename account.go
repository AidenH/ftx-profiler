package main

import (
	"net/http"
)

type AccountState struct {
	Result struct {
		PosSize float64 `json:"totalPositionSize"`
		Balance float64 `json:"totalAccountValue"`
	}
}

var Account = new(AccountState)
var client = &http.Client{}

// GetAccountInfo() retrieves user's FTX account info via http API request.
// Function will return an *http.Response for testing purposes to be used
// in conjunction with ParseHttpResp()
func GetAccountInfo() (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://ftx.com/api/account", nil)
	if err != nil {
		return nil, err
	}

	signature, ts, err := CreateSignature("/api/account")
	if err != nil {
		return nil, err
	}

	req.Header.Add("FTX-KEY", Api)
	req.Header.Add("FTX-SIGN", signature)
	req.Header.Add("FTX-TS", ts)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	Account.ParseHttpResp(resp)

	return resp, nil
}

func GetOpenPositions() (*http.Response, error) {
	return nil, nil
}
