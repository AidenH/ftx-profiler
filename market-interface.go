package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/sacOO7/gowebsocket"
)

type Request struct {
	Args    Args     `json:"args"`
	Params  []string `json:"params"`
	Method  string   `json:"method"`
	Op      string   `json:"op"`
	Channel string   `json:"channel"`
	Market  string   `json:"market"`
	Id      int      `json:"id"`
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

type OrdersResponse struct {
	Channel string    `json:"channel"`
	Type    string    `json:"type"`
	Data    ReplyData `json:"data"`
}

type ReplyData struct {
	Id            int     `json:"id"`
	Market        string  `json:"market"`
	Type          string  `json:"type"`
	Price         float64 `json:"price"`
	Size          float64 `json:"size"`
	Side          string  `json:"side"`
	Status        string  `json:"status"`
	FilledSize    float64 `json:"filledSize"`
	RemainingSize float64 `json:"remainingSize"`
	Liquidation   bool    `json:"liquidation"`
	AvgFillPrice  float64 `json:"avgFillPrice"`
	Time          string  `json:"time"`
	CreatedAt     string  `json:"createdAt"`
}

func SocketInit() error {

	var t TradesResponse
	var o OrdersResponse

	var sockErr error

	for _, i := range State.Connections {
		i.Socket.OnConnected = func(socket gowebsocket.Socket) {
			i.Subscribe(i)
		}

		i.Socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
			sockErr = err
		}

		i.Socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {

			if strings.Contains(msg, "orders") {
				_, sockErr = handleOrderReplies(o, msg)
				if sockErr != nil {
					FileWrite(sockErr.Error())
				}

			} else if strings.Contains(msg, "trades") {
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

			} else if strings.Contains(msg, "fills") {
				//FileWrite(fmt.Sprintln("fill:\n", msg))

			} else if !strings.Contains(msg, "pong") {
				sockErr = errors.New(fmt.Sprintf("unknown event type:\n %s\n", msg))
				FileWrite(fmt.Sprintln("unknown event type:\n", msg))

			}
		}

		i.Socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
			log.Println("Disconnected")
			sockErr = err
			return
		}

		i.Socket.Connect()
	}

	return sockErr

}

func handleOrderReplies(o OrdersResponse, msg string) (OrdersResponse, error) {
	// filter json data into OrdersResponse type
	if err := json.Unmarshal([]byte(msg), &o); err != nil {
		return OrdersResponse{}, err
	}

	if o.Data.Status == "new" {
		// handle new orders
		Account.Orders[o.Data.Id] = Order{
			Id:            o.Data.Id,
			Type:          o.Data.Type,
			Price:         o.Data.Price,
			Size:          o.Data.Size,
			Side:          o.Data.Side,
			FilledSize:    o.Data.FilledSize,
			RemainingSize: o.Data.RemainingSize,
		}

		FileWrite(fmt.Sprintln(Account.Orders))

	} else if o.Data.Status == "closed" {
		// handle closed orders
		delete(Account.Orders, o.Data.Id)
	}

	if err := PrintOrders(); err != nil {
		return OrdersResponse{}, err
	}

	return o, nil
}

func handleTradeReplies(t TradesResponse, msg string) error {

	// if SetMiddle is armed and last price is not init-0, add new last price
	//  middle-of-profile price
	if CState.SetMiddle && State.LastPrice != 0 {
		CState.Middle = State.LastPrice
		CState.SetMiddle = false
	}

	// function here to handle replies based on "type" field
	if err := json.Unmarshal([]byte(msg), &t); err != nil {
		return err
	}

	// if more than one Data element, iterate through items
	// else, just display one Data item
	// color code based on transaction's market side
	if len(t.Data) > 1 && t.Data[0].Size > State.VolMinFilter {

		// init variables for counting single cumulative trades across
		// 	multiple price levels
		cumulSize := 0.0
		var cumulSide string
		cumulPrice := 0.0

		for _, v := range t.Data {
			// if size is above user-configured minimum volume filter
			if v.Size > State.VolMinFilter {
				p, err := Round(v.Price, State.PricePrecision)
				if err != nil {
					return err
				}

				c := fmt.Sprintf("%.1f", v.Size)

				// add event data to VData var
				AddVData(p, v.Size)

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

	} else if len(t.Data) != 0 && t.Data[0].Size > State.VolMinFilter {

		p, err := Round(t.Data[0].Price, State.PricePrecision)
		if err != nil {
			return err
		}

		c := fmt.Sprintf("%.1f", t.Data[0].Size)

		// add event data to VData var
		AddVData(p, t.Data[0].Size)

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
