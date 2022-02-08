package main

import (
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
		fmt.Println(msg)
	}

	socket.Connect()

	for resp == "" {
		if resp != "" {
			fmt.Println(resp)
		}
	}
}

func TestOrdersSocket(t *testing.T) {
	State.Market = "GALA-PERP"

	socket := gowebsocket.New(SocketEndpoint)
	var resp string

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
		fmt.Println(msg)
	}

	socket.Connect()

	for resp == "" {
		if resp != "" {
			fmt.Println(resp)
		}
	}

}
