package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// bitmexChan := make(chan []BitmexData)
	mux := &Mux{
		operations: make(chan func(map[int]*Connection)),
	}
	go SubscribeToBitMex(mux)
	go mux.ManageConnections()
	startServer(mux)
}

func startServer(mux *Mux) {
	r := gin.Default()
	r.Use(errorHandler())
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, to use websocket connect to /ws via websocket connection")
	})
	r.GET("/ws", func(c *gin.Context) {
		WShandler(c, mux)
	})
	r.Run("localhost:8081")
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
}
