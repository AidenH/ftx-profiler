package main

import "github.com/sacOO7/gowebsocket"

const BinSocketEndpoint = "wss://fstream.binance.com/"

func NewBinConnection() Connection {
	return Connection{
		Name:      "binance",
		Socket:    gowebsocket.New(BinSocketEndpoint),
		Subscribe: BinSubscribe,
	}
}

func BinSubscribe(c Connection) {

}
