package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dush-t/goloca/db"
	"github.com/dush-t/goloca/pool"
)

// AddDataPoint will insert the sent datapoint into the database
func AddDataPoint(p pool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var reqBody db.RideDataPoint

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			payload := failMessage("invalid_data_format")
			json.NewEncoder(w).Encode(payload)
			return
		}

		err = reqBody.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			payload := failMessage(fmt.Sprint(err))
			json.NewEncoder(w).Encode(payload)
			return
		}

		job := pool.CreateJob(reqBody)

		p.Dispatch(job)
		log.Println("Job dispatched")
		w.Header().Set("Content-Type", "application/json")

		status := <-job.StatusChan

		var payload responseMessage
		s, isBool := status.(bool)
		err, isErr := status.(error)

		if isBool {
			if s {
				w.WriteHeader(http.StatusCreated)
				payload = successMessage(reqBody.ID)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				payload = failMessage("unknown_internal_server_error")
			}
		} else if isErr {
			w.WriteHeader(http.StatusInternalServerError)
			payload = failMessage(fmt.Sprintf("%s", err))
		}

		json.NewEncoder(w).Encode(payload)
	})
}
