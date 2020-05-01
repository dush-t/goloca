package main

import (
	"log"
	"net/http"

	"github.com/dush-t/goloca/api"
)

func main() {
	config := InitializeApp()

	http.Handle("/add_datapoint", api.AddDataPoint(config.Pool))

	log.Println("HTTP server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
