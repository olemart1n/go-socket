package main

import (
	"context"
	"log"
	"net/http"
)

func main () {
	setupAPI()

	// log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServeTLS(":8443", "example.com+4.pem", "example.com+4-key.pem", nil))
	
}

func setupAPI () {
	ctx:= context.Background()
	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/login", manager.loginHandler)
	http.HandleFunc("/ws", manager.serveWS)
	
}