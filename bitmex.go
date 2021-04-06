package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	API_KEY    = "ORqVaoVf1TJrVnKexpWjHfjk"
	API_SECRET = "mvK7p-zYF5He2eistXxXUvASoJWRGvp6eOO5TF2gn4BHI2iB"
)

type BitmexResponse struct {
	Data []BitmexData `json:"data"`
}

type BitmexData struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"markPrice"`
	TimeStamp string  `json:"timestamp"`
}

func SubscribeToBitMex(bitmexChan chan []BitmexData) {
	// Establishing connection to bitmex
	conn, _, err := websocket.DefaultDialer.Dial("wss://testnet.bitmex.com/realtime", nil)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	defer conn.Close()

	//Authorizing
	expires := int(time.Now().Unix() + 3600)
	auth := &BitmexCommand{
		Op:   "authKeyExpires",
		Args: []interface{}{API_KEY, expires, GetSignature(expires)},
	}
	conn.WriteJSON(auth)

	//Subscribing to the instrument
	sub := &BitmexCommand{
		Op:   "subscribe",
		Args: []interface{}{"instrument"},
	}
	conn.WriteJSON(sub)

	//Processing received messages
	for {
		response := &BitmexResponse{}
		// _, msg, err := conn.ReadMessage()
		err = conn.ReadJSON(response)
		if err != nil {
			log.Println("ERROR: ", err)
			continue
		}
		// log.Println(string(msg))
		// log.Println("response from bitmex: ", response)
		bitmexChan <- response.Data
	}
}

func GetSignature(expires int) string {
	h := hmac.New(sha256.New, []byte(API_SECRET))
	data := "GET/realtime" + strconv.Itoa(expires)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
