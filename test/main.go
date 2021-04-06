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
		Action  string `json:"action"`
		Symbols string `json:"symbols"`
	}

	conn.WriteJSON(&Sub{
		Action:  "subscribe",
		Symbols: "HASUSD",
	})
	subscription := &Sub{}
	conn.ReadJSON(subscription)
	fmt.Println(subscription)
	<-time.After(10 * time.Second)
}
