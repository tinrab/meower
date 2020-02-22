package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	id       int
	socket   *websocket.Conn
	outbound chan []byte
}

func newClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (c *Client) write() {
	for {
		select {
		case data, ok := <-c.outbound:
			if !ok {
				if err := c.socket.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Print(err)
				}
				return
			}
			if err := c.socket.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Print(err)
				c.hub.disconnect(c)
			}
		}
	}
}

func (c Client) close() {
	if err := c.socket.Close(); err != nil {
		log.Fatal(err)
	}
	close(c.outbound)
}
