package main

import (
	"github.com/gin-gonic/gin"
)

const (
	API_KEY    = "ORqVaoVf1TJrVnKexpWjHfjk"
	API_SECRET = "mvK7p-zYF5He2eistXxXUvASoJWRGvp6eOO5TF2gn4BHI2iB"
)

func main() {
	// SubscribeToBitMex()
	startServer()
}

func startServer() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, to use websocket connect to /ws via websocket connection")
	})
	r.GET("/ws", func(c *gin.Context) {
		WShandler(c.Writer, c.Request)
	})
	r.Run("localhost:8081")
}
