package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sacOO7/gowebsocket"
)

const BinSocketBase = "wss://fstream.binance.com/"

func NewBinConnection() Connection {
	return Connection{
		Name:      "binance",
		Socket:    gowebsocket.New(BinSocketBase),
		Subscribe: BinSubscribe,
	}
}

func BinSubscribe(c Connection) error {
	market := strings.Replace(State.Market, "-", "", -1)

	log.Println("connected to binance websocket")

	dat, err := json.Marshal(Request{
		Method: "SUBSCRIBE",
		Params: []string{fmt.Sprint(market, "@aggTrade")},
	})
	if err != nil {
		return err
	}

	c.Socket.SendBinary(dat)

	return nil
}
