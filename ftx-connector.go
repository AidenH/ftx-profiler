package main

import (
	"encoding/json"
	"log"

	"github.com/sacOO7/gowebsocket"
)

const FTXSocketEndpoint = "wss://ftx.com/ws/"

func NewFTXConnection() Connection {
	return Connection{
		Name:      "ftx",
		Socket:    gowebsocket.New(FTXSocketEndpoint),
		Subscribe: FTXSubscribe,
	}
}

func FTXSubscribe(c Connection) {
	log.Println("Connected!")

	// send json ping to server
	FTXPingRequest(c.Socket)

	FTXAuthStreamLogin(c.Socket)

	// attempt sub to trades, fills, orders
	FTXSubscribeRequest(c.Socket, "trades")
	FTXSubscribeRequest(c.Socket, "fills")
	FTXSubscribeRequest(c.Socket, "orders")
}

func FTXPingRequest(s gowebsocket.Socket) error {

	dat, err := json.Marshal(Request{
		Op: "ping",
	})
	if err != nil {
		return err
	}

	s.SendBinary(dat)

	return nil
}

func FTXAuthStreamLogin(s gowebsocket.Socket) error {
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
func FTXSubscribeRequest(s gowebsocket.Socket, typ string) error {
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

	s.SendBinary(dat)

	return nil
}

func FTXUnsubscribeRequest(s gowebsocket.Socket, ch string) error {
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
