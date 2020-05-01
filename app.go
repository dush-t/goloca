package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/dush-t/goloca/db"
	"github.com/dush-t/goloca/pool"
)

// Config stores global variables that the app needs. These
// variables will be passed to functions using dependency
// injection
type Config struct {
	Store db.Store
	Pool  pool.Pool
}

func baseDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func buildDBInsertAction(conf Config) pool.WorkerAction {
	return func(job pool.Job) error {
		err := conf.Store.Insert(job.DataPoint)
		if err != nil {
			return err
		}
		job.CompleteChan <- true
		return nil
	}
}

// InitializeApp will create a database connection and start the threadpool
func InitializeApp() Config {

	// Load environment variables
	baseDir := baseDir()
	LoadEnv(filepath.Join(baseDir, "config.env"))
	var conf Config

	// Connect to datastore.
	store := db.CreateStore()
	conf.Store = store

	// Create worker pool
	var p pool.Pool
	numWorkers, _ := strconv.Atoi(os.Getenv("NUM_POOL_WORKERS"))
	p.Start(numWorkers, buildDBInsertAction(conf))

	conf.Pool = p

	return conf
}
