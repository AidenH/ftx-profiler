package main

import (
	"fmt"
	"net/http"
)

type Order = map[string]interface{}

type Orders struct {
	Success bool
	Result  []struct {
		Id            int
		Type          string
		Side          string
		Price         float64
		Size          float64
		RemainingSize float64
	}
}

type AccountState struct {
	Success bool
	Error   string
	Result  struct {
		PosSize float64 `json:"totalPositionSize"`
		Balance float64 `json:"totalAccountValue"`

		// position data returned from http request
		PositionsData []struct {
			Future string  `json:"future"`
			Entry  float64 `json:"entryPrice"`
			Size   float64 `json:"openSize"`
			Side   string  `json:"side"`
			Pnl    float64 `json:"unrealizedPnl"`
		}
	}

	// open position populated during RetrieveAccountInfo
	Position struct {
		Entry float64
		Size  float64
		Side  string
		Pnl   float64
	}

	// open orders
	Orders []struct {
		Price float64
		Size  float64
		Side  string
	}
}

// GetAccountInfo() retrieves user's FTX account info via http API request.
// Function will return an *http.Response for testing purposes to be used
// in conjunction with ParseHttpResp()
func (a *AccountState) GetAccountInfo() error {
	req, err := http.NewRequest("GET", "https://ftx.com/api/account", nil)
	if err != nil {
		return err
	}

	signature, ts, err := CreateHttpSignature("/api/account")
	if err != nil {
		return err
	}

	// add headers
	req.Header.Add("FTX-KEY", Api)
	req.Header.Add("FTX-SIGN", signature)
	req.Header.Add("FTX-TS", ts)

	// send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// parse request response into JSON
	ParseHttpResp(resp, a)

	return nil
}

// GetOpenOrders() returns a user's open resting orders
func (o *Orders) GetOpenOrders() error {
	url := fmt.Sprintf("https://ftx.com/api/orders?market=%s", "GALA-PERP")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	signature, ts, err := CreateHttpSignature(
		fmt.Sprintf("/api/orders?market=%s", "GALA-PERP"))
	if err != nil {
		return err
	}

	// add headers
	req.Header.Add("FTX-KEY", Api)
	req.Header.Add("FTX-SIGN", signature)
	req.Header.Add("FTX-TS", ts)

	// send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// parse request response into JSON
	if err := ParseHttpResp(resp, o); err != nil {
		return err
	}

	return nil
}
