package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var WSupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Main websocket handler for clients
func WShandler(c *gin.Context, mux *Mux) {
	conn, err := WSupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(fmt.Errorf("Failed to set websocket upgrade: %v.\n", err))
		return
	}
	defer conn.Close()

	subChan := readConnection(conn)
	connection := &Connection{
		Info: make(chan Info),
		// SubSymbols:   make(chan []string),
		// UnsubSymbols: make(chan []string),
		Symbols: make(map[string]struct{}),
		Id:      int(time.Now().UnixNano()), //not ideal but kinda unique for the matter
	}
	for {
		select {
		case subscription, ok := <-subChan:
			if !ok {
				mux.RemoveConnection(connection)
				return
			}
			processSub(mux, subscription, connection, conn)
		case info := <-connection.Info:
			err := conn.WriteJSON(info)
			if err != nil {
				mux.RemoveConnection(connection)
				return
			}
		}

	}

}

//reads all incoming jsons from client connection
func readConnection(conn *websocket.Conn) chan Sub {
	subChan := make(chan Sub)
	go func() {
		for {
			subscription := &Sub{}
			err := conn.ReadJSON(subscription)
			if err != nil {
				log.Println("ERROR: ", err)
				close(subChan)
				break
			}
			subChan <- *subscription
		}
	}()
	return subChan
}

//distributes logic by action
func processSub(mux *Mux, subscription Sub, conn *Connection, webConn *websocket.Conn) {
	switch subscription.Action {
	case "subscribe":
		mux.AddConnection(conn)
		mux.Subscribe(conn.Id, subscription.Symbols)
	case "unsubscribe":
		if len(subscription.Symbols) == 0 {
			mux.RemoveConnection(conn)
			return
		}
		mux.Unsubscribe(conn.Id, subscription.Symbols)
	default:
		webConn.WriteJSON(Error{Message: "Unknown command was sent: " + subscription.Action})
		log.Println("Unknown Command: ", subscription.Action)
	}
}
