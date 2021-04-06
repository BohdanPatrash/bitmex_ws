package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func SubscribeToBitMex() {
	// Establishing connection to bitmex
	conn, _, err := websocket.DefaultDialer.Dial("wss://testnet.bitmex.com/realtime", nil)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	defer conn.Close()

	//Authorizing
	expires := int(time.Now().Unix() + 3600)
	auth := &Command{
		Op:   "authKeyExpires",
		Args: []interface{}{API_KEY, expires, GetSignature(expires)},
	}
	conn.WriteJSON(auth)

	//Subscribing to the instrument
	sub := &Command{
		Op:   "subscribe",
		Args: []interface{}{"instrument"},
	}
	conn.WriteJSON(sub)

	//Processing received messages
	for {
		<-time.After(2 * time.Second)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR: ", err)
			break
		}
		fmt.Println("recieved <- : ", string(msg))
	}
}

func GetSignature(expires int) string {
	h := hmac.New(sha256.New, []byte(API_SECRET))
	data := "GET/realtime" + strconv.Itoa(expires)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
