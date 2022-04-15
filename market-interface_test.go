package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sacOO7/gowebsocket"
)

var tinkering = true

func TestSocketInit(t *testing.T) {
	if err := SocketInit(); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestFTXTradesSocket(t *testing.T) {
	State.Market = "BTC-PERP"
	State.Connections = append(State.Connections, NewFTXConnection())

	c := State.Connections[0]
	var resp string
	var tr TradesResponse

	c.Socket.OnConnected = func(socket gowebsocket.Socket) {
		c.Subscribe(c)
	}

	c.Socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println(err)
		t.Fail()
	}

	c.Socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
		json.Unmarshal([]byte(msg), &tr)
		fmt.Println(tr)
	}

	c.Socket.Connect()

	for resp == "" {
		// pass on successful subscribe
		if tr.Type == "subscribed" && !tinkering {
			return
		}
	}
}

func TestFTXOrdersSocket(t *testing.T) {
	State.Market = "BTC-PERP"
	State.Connections = append(State.Connections, NewFTXConnection())

	c := State.Connections[0]
	var resp string
	var ord OrdersResponse
	var err error

	c.Socket.OnConnected = func(socket gowebsocket.Socket) {
		if err = FTXAuthStreamLogin(socket); err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}

		c.Subscribe(c)
	}

	c.Socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println(err)
		t.Fail()
	}

	c.Socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
		fmt.Println(msg)

		ord, err = handleOrderReplies(ord, msg)
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
	}

	c.Socket.Connect()

	for resp == "" {
		// pass on successful subscribe
		if ord.Type == "subscribed" && !tinkering {
			return
		}
	}

}
