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
	Args    Args   `json:"args"`
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type Args struct {
	Key  string `json:"key"`
	Sign string `json:"sign"`
	Time int64  `json:"time"`
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

func SocketInit() error {

	socket := gowebsocket.New(SocketEndpoint)
	var t TradesResponse

	var sockErr error

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected!")

		// send json ping to server
		pingRequest(socket)

		// attempt sub to trades, fills, orders
		subscribeRequest(socket, "trades")
		subscribeRequest(socket, "fills")
		subscribeRequest(socket, "orders")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		sockErr = err
	}

	socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
		if !CState.LockWrite {
			CState.LockWrite = true

			go func() {
				sockErr = handleTradeReplies(t, msg)
				if sockErr != nil {
					FileWrite(sockErr.Error())
				}
				if State.ProfileTrue {
					PrintProfile()
				}
				SetStatus()

				CState.LockWrite = false
			}()
		}
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected")
		sockErr = err
		return
	}

	socket.Connect()

	return sockErr

}

func handleTradeReplies(t TradesResponse, msg string) error {

	// if SetMiddle is armed and last price is not init-0, add new last price
	//  middle-of-profile price
	if CState.SetMiddle && State.LastPrice != 0 {
		CState.Middle = State.LastPrice
		CState.SetMiddle = false
	}

	// function here to handle replies based on "type" field
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

			p, err := Round(v.Price, State.PricePrecision)
			if err != nil {
				return err
			}

			c := fmt.Sprintf("%.1f", v.Size)

			// add event data to VData var
			AddVData(p, v.Size)
			//FileWrite(fmt.Sprintf("VData: %f - %f", p, v.Size))

			// if side buy, print green
			if v.Side == "buy" {
				if !State.Aggregate {
					PrintTape("buy", p, c)

				} else {
					cumulSize += v.Size
					cumulSide = "buy"
					cumulPrice = v.Price
				}

			} else if v.Side == "sell" {
				// if side sell, print red
				if !State.Aggregate {
					PrintTape("sell", p, c)

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
		if State.Aggregate {

			p, err := Round(cumulPrice, State.PricePrecision)
			if err != nil {
				return err
			}

			c := fmt.Sprintf("%.1f", cumulSize)
			if cumulSide == "buy" {
				PrintTape("buy", p, c)

			} else if cumulSide == "sell" {
				PrintTape("sell", p, c)

			}

			// set session Open price if last price was init 0
			if State.LastPrice == 0 {
				State.OpenPrice = p
			}

			State.LastPrice = p
		}

	} else if len(t.Data) != 0 {

		p, err := Round(t.Data[0].Price, State.PricePrecision)
		if err != nil {
			return err
		}

		c := fmt.Sprintf("%.1f", t.Data[0].Size)

		// add event data to VData var
		AddVData(p, t.Data[0].Size)
		//FileWrite(fmt.Sprintf("VData: %f - %f", p, t.Data[0].Size))

		// set session Open price if last price was init 0
		if State.LastPrice == 0 {
			State.OpenPrice = p
		}

		// set global last price
		State.LastPrice = p

		if t.Data[0].Side == "buy" {
			PrintTape("buy", p, c)

		} else if t.Data[0].Side == "sell" {
			PrintTape("sell", p, c)

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

func AuthStreamLogin(s gowebsocket.Socket) error {
	signature, ts, err := CreateSocketSignature()
	if err != nil {
		return err
	}

	// args for authorization
	Auth := Args{
		Key:  Api,
		Sign: signature,
		Time: ts,
	}

	// rest of packet including args
	dat, err := json.Marshal(Request{
		Args: Auth,
		Op:   "login",
	})
	if err != nil {
		return err
	}

	s.SendBinary(dat)

	return nil
}

// subscribeRequest connects websocket client to provided FTX stream
func subscribeRequest(s gowebsocket.Socket, typ string) error {
	var Auth Args

	// trades subscribe packet
	dat, err := json.Marshal(Request{
		Args:    Auth,
		Op:      "subscribe",
		Channel: typ,
		Market:  State.Market,
	})
	if err != nil {
		return err
	}

	//fmt.Println(string(dat))

	s.SendBinary(dat)

	return nil
}

func unsubscribeRequest(s gowebsocket.Socket, ch string) error {
	dat, err := json.Marshal(Request{
		Op:      "unsubscribe",
		Channel: ch,
		Market:  State.Market,
	})
	if err != nil {
		return err
	}

	s.SendBinary(dat)

	return nil
}
