package main

import (
	"log"
	"net/http"

	"github.com/dush-t/goloca/api"
)

func main() {
	config := InitializeApp()

	http.Handle("/add_datapoint", api.AddDataPoint(config.Pool))

	// WebSocket, to be used during peak hours to save the overhead of establishing
	// a new connection everytime
	http.HandleFunc("/ws", api.AddDataPointSocket(config.SocketUpgrader, config.Pool))

	log.Println("HTTP server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
