package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	socket *websocket.Conn
	send   chan *message
	room   *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	err := c.socket.Close()
	if err != nil {
		return 
	}
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	err := c.socket.Close()
	if err != nil {
		return 
	}
}
