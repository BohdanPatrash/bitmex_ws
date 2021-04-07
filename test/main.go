//just a file for manual testing
package main

import (
	"fmt"

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
		Error     string  `jsno:"error"`
	}
	//subing to instrument:XBTUSD
	conn.WriteJSON(&Sub{
		Action:  "subscribe",
		Symbols: []string{"XBTUSD"},
	})

	for i := 0; i < 5; i++ {
		info := &Info{}
		conn.ReadJSON(info)
		fmt.Println(info)
	}

	//wrong sub command
	conn.WriteJSON(&Sub{
		Action:  "asd",
		Symbols: []string{},
	})

	info := &Info{}
	conn.ReadJSON(info)
	fmt.Println(info)

	//unsubing
	conn.WriteJSON(&Sub{
		Action:  "unsubscribe",
		Symbols: []string{},
	})
}
