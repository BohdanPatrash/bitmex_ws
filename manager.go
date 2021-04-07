package main

import (
	"log"
	"sync"
)

type Clients struct {
	mu          sync.Mutex
	connections map[int]*Connection
}

type Connection struct {
	Id           int
	Symbols      map[string]struct{}
	Info         chan Info
	UnsubSymbols chan []string
	SubSymbols   chan []string
}

//clients that are connected
var clients *Clients = &Clients{
	connections: make(map[int]*Connection),
}

//ManageConnections sends relative data to the subscribers
func ManageConnections(bitmexChan chan []BitmexData) {
	for {
		data := <-bitmexChan
		clients.mu.Lock()

		for _, conn := range clients.connections {
			for _, val := range data {
				if _, ok := conn.Symbols[val.Symbol]; ok || len(conn.Symbols) == 0 {
					conn.Info <- Info{
						Price:     val.Price,
						Symbol:    val.Symbol,
						Timestamp: val.TimeStamp,
					}
				}
			}

		}

		clients.mu.Unlock()
	}
}

//AddConnection adds connection if it doesn't exist already
func AddConnection(c *Connection) {
	clients.mu.Lock()
	defer clients.mu.Unlock()
	if _, ok := clients.connections[c.Id]; ok {
		return
	}
	clients.connections[c.Id] = c
	listenSubUpdate(clients.connections[c.Id])
	log.Printf("connection ID:%v successfully subscribed!", c.Id)
}

//RemoveConnection removes existing connection
func RemoveConnection(c *Connection) {
	clients.mu.Lock()
	delete(clients.connections, c.Id)
	log.Printf("connection ID:%v successfully unsubscribed!", c.Id)
	clients.mu.Unlock()
}

//listenSubUpdate listening to updates on existing connections and updates subscribed symbols
func listenSubUpdate(c *Connection) {
	go func() {
		for {
			select {
			case symbols := <-c.SubSymbols:
				clients.mu.Lock()
				for _, v := range symbols {
					c.Symbols[v] = struct{}{}
					log.Printf("connection ID:%v successfully subscribed to %v", c.Id, v)
				}
				clients.mu.Unlock()
			case symbols := <-c.UnsubSymbols:
				clients.mu.Lock()
				for _, v := range symbols {
					delete(c.Symbols, v)
					log.Printf("connection ID:%v successfully unsubscribed from %v", c.Id, v)
				}
				clients.mu.Unlock()
			}
		}
	}()
}
