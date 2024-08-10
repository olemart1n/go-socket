package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	connection *websocket.Conn
	manager *Manager
	egress chan Event
}

type ClientList map[*Client]bool 

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager: manager,
		egress: make(chan Event),
	}
}

func (c *Client) readMessages () {
	defer func() {
		//cleanup connection function that is deffered (happens when the loop breaks)
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
	}

	c.connection.SetReadLimit(512)
	c.connection.SetPongHandler(c.pongHandler)
	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event :%v", err)
			break
		}
		
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error when handling message: ", err)
		}
	}
}
func (c *Client) writeMessages () {
	defer func() {
		//cleanup connection function that is deffered (happens when the loop breaks)
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)


	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed: ",  err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err:= c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")

		case <- ticker.C:
			log.Println("ping")

			// send ping to the client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("writeMessage err: ", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler (pongMessage string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}