package main

import (
	// "fmt"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)
var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {

		// 	origin := r.Header.Get("Origin")
		// 	switch origin {
		// 	case "http://127.0.0.1:5500":
				
		// 	}
			return true
		},
		
	}

)

type Manager struct {
	clients ClientList
	sync.RWMutex
	otps RetentionMap
	handlers map[string]EventHandler
}

func (m *Manager) loginHandler(w http.ResponseWriter, r *http.Request) {
	type userLoginRequest struct {
		Username string `json: "username"` 
		Password string `json: "password"`
	}
	

	var req userLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
	if req.Username == "percy" && req.Password == "123" {
		type response struct {
			OTP string `json: "otp"`
		}
		otp := m.otps.NewOTP()

		res := response{
			OTP: otp.Key,
		}
		data, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	w.WriteHeader(401)
	// w.WriteHeader(http.StatusUnauthorized)
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
		otps: NewRetentionMap(ctx, 5*time.Second),
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

	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(r.URL.RawQuery)
		return
	}

	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		
		return
	}
	log.Println("new connection")

	// if the code below executes, the user has a valid ticket(credentials)
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

