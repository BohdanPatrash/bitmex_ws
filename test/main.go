package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	conn, res, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	defer conn.Close()
	defer res.Body.Close()

	type Sub struct {
		Action  string   `json:"action"`
		Symbols []string `json:"symbols"`
	}

	type Info struct {
		Timestamp string  `json:"timestamp"`
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"price"`
	}

	conn.WriteJSON(&Sub{
		Action:  "subscribe",
		Symbols: []string{},
	})

	for i := 0; i < 200; i++ {
		info := &Info{}
		conn.ReadJSON(info)
		fmt.Println(info)
	}
	<-time.After(2 * time.Second)
	conn.WriteJSON(&Sub{
		Action:  "unsubscribe",
		Symbols: []string{},
	})
}
