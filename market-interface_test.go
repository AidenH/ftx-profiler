package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/sacOO7/gowebsocket"
)

func TestSocketInit(t *testing.T) {
	if err := SocketInit(); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestTradesSocket(t *testing.T) {
	State.Market = "BTC-PERP"

	socket := gowebsocket.New(SocketEndpoint)
	var resp string
	var tr TradesResponse

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("connected!")
		pingRequest(socket)

		if err := subscribeRequest(socket, "trades"); err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println(err)
		t.Fail()
	}

	socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
		json.Unmarshal([]byte(msg), &tr)
		fmt.Println(tr)
	}

	socket.Connect()

	for resp == "" {
		if tr.Type == "subscribed" {
			return
		}
	}
}

func TestOrdersSocket(t *testing.T) {
	State.Market = "BTC-PERP"

	socket := gowebsocket.New(SocketEndpoint)
	var resp string
	var ord OrdersResponse

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("connected!")
		pingRequest(socket)

		if err := AuthStreamLogin(socket); err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}

		if err := subscribeRequest(socket, "orders"); err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println(err)
		t.Fail()
	}

	socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
		json.Unmarshal([]byte(msg), &ord)
		fmt.Println(ord)
	}

	socket.Connect()

	for resp == "" {
		if ord.Type == "subscribed" {
			return
		}
	}

}
