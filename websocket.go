package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var WSupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WShandler(w http.ResponseWriter, r *http.Request) {
	conn, err := WSupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %v.\n", err)
		return
	}

	for {
		<-time.After(2 * time.Second)
		subscription := &Sub{}
		err := conn.ReadJSON(subscription)
		if err != nil {
			fmt.Println("ERROR: ", err)
			break
		}
		conn.WriteJSON(subscription)
	}
}
