package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dush-t/goloca/db"
	"github.com/dush-t/goloca/pool"
)

// AddDataPoint will insert the sent datapoint into the database
func AddDataPoint(p pool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody db.RideDataPoint

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		completeChan := make(chan bool)
		job := pool.Job{CompleteChan: completeChan, DataPoint: reqBody}

		p.Dispatch(job)
		log.Println("Job dispatched")
		w.Header().Set("Content-Type", "application/json")

		complete := <-completeChan
		if complete {
			w.WriteHeader(http.StatusCreated)
			payload := struct {
				ID string `json:"id"`
			}{ID: reqBody.ID}
			json.NewEncoder(w).Encode(payload)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
