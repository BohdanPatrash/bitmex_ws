package main

import (
	"log"
)

type Connection struct {
	Id           int
	Symbols      map[string]struct{}
	Info         chan Info
	UnsubSymbols chan []string
	SubSymbols   chan []string
}

type Mux struct {
	operations chan func(map[int]*Connection)
}

func (m *Mux) ManageConnections() {
	connections := make(map[int]*Connection)
	for operation := range m.operations {
		operation(connections)
	}
}

func (m *Mux) AddConnection(c *Connection) {
	m.operations <- func(conns map[int]*Connection) {
		if _, ok := conns[c.Id]; ok {
			return
		}
		conns[c.Id] = c
		log.Printf("connection ID:%v successfully subscribed!", c.Id)
	}
}

func (m *Mux) RemoveConnection(c *Connection) {
	m.operations <- func(conns map[int]*Connection) {
		delete(conns, c.Id)
		log.Printf("connection ID:%v successfully unsubscribed!", c.Id)
	}
}

func (m *Mux) BitmexUpdate(data []BitmexData) {
	m.operations <- func(conns map[int]*Connection) {
		for _, conn := range conns {
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
	}
}

func (m *Mux) Subscribe(id int, symbols []string) {
	m.operations <- func(conns map[int]*Connection) {
		for _, v := range symbols {
			conns[id].Symbols[v] = struct{}{}
			log.Printf("connection ID:%v successfully subscribed to %v", id, v)
		}
	}
}

func (m *Mux) Unsubscribe(id int, symbols []string) {
	m.operations <- func(conns map[int]*Connection) {
		for _, v := range symbols {
			delete(conns[id].Symbols, v)
			log.Printf("connection ID:%v successfully unsubscribed from %v", id, v)
		}
	}
}
