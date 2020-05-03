package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/dush-t/goloca/db"
	"github.com/dush-t/goloca/pool"
)

// Config stores global variables that the app needs. These
// variables will be passed to functions using dependency
// injection
type Config struct {
	Store          db.Store
	Pool           pool.Pool
	SocketUpgrader websocket.Upgrader
}

func buildDBInsertAction(conf Config) pool.WorkerAction {
	return func(job pool.Job) error {
		err := conf.Store.Insert(job.DataPoint)
		if err != nil {
			return err
		}
		job.StatusChan <- true
		return nil
	}
}

// InitializeApp will create a database connection and start the threadpool
func InitializeApp() Config {

	// Load environment variables
	LoadEnv("config.env")
	var conf Config

	// Connect to datastore.
	store := db.CreateStore()
	conf.Store = store

	// Create worker pool
	var p pool.Pool
	numWorkers, _ := strconv.Atoi(os.Getenv("NUM_POOL_WORKERS"))
	p.Start(numWorkers, buildDBInsertAction(conf))

	conf.Pool = p

	// Setup websocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conf.SocketUpgrader = upgrader

	return conf
}
