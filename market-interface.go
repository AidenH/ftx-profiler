package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/sacOO7/gowebsocket"
)

const SocketEndpoint = "wss://ftx.com/ws/"

type Request struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type TradesResponse struct {
	Channel string      `json:"channel"`
	Market  string      `json:"market"`
	Type    string      `json:"type"`
	Data    []ReplyData `json:"data"`
}

type ReplyData struct {
	Id          int     `json:"id"`
	Price       float64 `json:"price"`
	Size        float64 `json:"size"`
	Side        string  `json:"side"`
	Liquidation bool    `json:"liquidation"`
	Time        string  `json:"time"`
}

func SocketInit(state ProfileState) error {

	socket := gowebsocket.New(SocketEndpoint)
	var t TradesResponse

	var sockErr error

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected!")

		// send json ping to server
		pingRequest(socket)
		subscribeRequest(socket, state.Market)

		//clear terminal
		fmt.Println("\033[H\033[2J")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		sockErr = err
	}

	socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
		handleTradeReplies(t, msg, state)
		PrintProfile(state.Gui)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected")
		sockErr = err
		return
	}

	socket.Connect()

	return sockErr

}

func handleTradeReplies(t TradesResponse, msg string, state ProfileState) error {
	prec := state.PricePrecision

	//function here to handle replies based on "type" field
	json.Unmarshal([]byte(msg), &t)

	// if more than one Data element, iterate through items
	// else, just display one Data item
	// color code based on transaction's market side
	if len(t.Data) > 1 {

		// init variables for counting single cumulative trades across
		// 	multiple price levels
		cumulSize := 0.0
		var cumulSide string
		cumulPrice := 0.0

		for _, v := range t.Data {

			p := Round(v.Price, state.PricePrecision)
			c := fmt.Sprintf("%.1f", v.Size)

			VData[p] += v.Size

			// if side buy, print green
			if v.Side == "buy" {
				if !state.Aggregate {
					PrintTape(state.Gui, "buy", p, c, prec)

				} else {
					cumulSize += v.Size
					cumulSide = "buy"
					cumulPrice = v.Price
				}

			} else if v.Side == "sell" {
				// if side sell, print red
				if !state.Aggregate {
					PrintTape(state.Gui, "sell", p, c, prec)

				} else {
					cumulSize += v.Size
					cumulSide = "sell"
					cumulPrice = v.Price
				}
			} else {
				err := errors.New("handleTradeReplies - invalid side type")
				return err
			}
		}

		// if aggregation turned on, output cumulative event volume after
		// 	counting loop is complete
		if state.Aggregate {

			p := Round(cumulPrice, state.PricePrecision)
			c := fmt.Sprintf("%.1f", cumulSize)

			if cumulSide == "buy" {
				PrintTape(state.Gui, "buy", p, c, prec)

			} else if cumulSide == "sell" {
				PrintTape(state.Gui, "sell", p, c, prec)

			}
		}

	} else if len(t.Data) != 0 {

		p := Round(t.Data[0].Price, 3)
		c := fmt.Sprintf("%.1f", t.Data[0].Size)

		if t.Data[0].Side == "buy" {
			PrintTape(state.Gui, "buy", p, c, prec)

		} else if t.Data[0].Side == "sell" {
			PrintTape(state.Gui, "sell", p, c, prec)

		} else {
			err := errors.New("handleTradeReplies 1 item Data - invalid side type")
			return err
		}
	}

	return nil

}

func pingRequest(s gowebsocket.Socket) error {

	dat, err := json.Marshal(Request{
		Op: "ping",
	})
	if err != nil {
		return err
	}

	s.SendBinary(dat)

	return nil
}

func subscribeRequest(s gowebsocket.Socket, market string) error {
	dat, err := json.Marshal(Request{
		Op:      "subscribe",
		Channel: "trades",
		Market:  market,
	})
	if err != nil {
		return err
	}

	s.SendBinary(dat)

	return nil
}

func unsubscribeRequest(s gowebsocket.Socket) error {
	dat, err := json.Marshal(Request{
		Op:      "unsubscribe",
		Channel: "trades",
		Market:  "NEAR-PERP",
	})
	if err != nil {
		return err
	}

	s.SendBinary(dat)

	return nil
}
