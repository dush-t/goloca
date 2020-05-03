package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dush-t/goloca/db"
	"github.com/dush-t/goloca/pool"

	"github.com/gorilla/websocket"
)

// AddDataPointSocket will open a new websocket connection to the
// the server. Clients can use this connection to push data to the server.
func AddDataPointSocket(u websocket.Upgrader, p pool.Pool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := u.Upgrade(w, r, nil)
		log.Println("New websocker connection initiated by")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		defer ws.Close()

		for {
			var msg db.RideDataPoint
			err := ws.ReadJSON(&msg)

			if err != nil {
				log.Printf("Client sent invalid data")
				response := failMessage("invalid_data_format")
				ws.WriteJSON(response)
				continue
			}

			valErr := msg.Validate()
			if valErr != nil {
				response := failMessage(fmt.Sprint(valErr))
				ws.WriteJSON(response)
				continue
			}

			job := pool.CreateJob(msg)

			p.Dispatch(job)
			log.Println("Job dispatched")

			status := <-job.StatusChan
			var payload responseMessage
			s, isBool := status.(bool)
			err, isErr := status.(error)

			if isBool {
				if s {
					payload = successMessage(msg.ID)
				} else {
					payload = failMessage("unknown_internal_server_error")
				}
			} else if isErr {
				payload = failMessage(fmt.Sprintf("%s", err))
			}

			ws.WriteJSON(payload)
		}
	}
}
