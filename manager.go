package main

import (
	// "fmt"
	"errors"
	"fmt"
	"log"
	"net/http"

	"sync"

	"github.com/gorilla/websocket"
)
var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

)

type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("There is no such event type")
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request){
	log.Println("new connection")
	// updgrade regular http connection into websocket
	 conn, err := websocketUpgrader.Upgrade(w, r, nil)

	 if err != nil {
		log.Println(err)
		return
	 }
	 
	client := NewClient(conn, m)
	m.addClient(client)

	// // start two processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(c *Client) {
	m.Lock()
	defer m.Unlock()
	fmt.Println("client added")
	m.clients[c] = true
	
}
func (m *Manager) removeClient(c *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
		fmt.Println("client removed")
	}
	
}

