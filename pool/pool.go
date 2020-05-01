package pool

import (
	"log"
	"os"
	"strconv"
)

// Pool represents a thread-pool like construct
type Pool struct {
	NumWorkers int
	JobChannel *(chan Job)
	Workers    []*worker
}

// Start will spawn a given number of workers with a given
// WorkerAction
func (p *Pool) Start(numWorkers int, action WorkerAction) {
	(*p).NumWorkers = numWorkers

	maxChannelLen, _ := strconv.Atoi(os.Getenv("MAX_JOB_QUEUE_LENGTH"))
	jobChannel := make(chan Job, maxChannelLen)
	(*p).JobChannel = &jobChannel

	for i := 0; i < numWorkers; i++ {
		stopChan := make(chan bool)
		w := worker{ID: i, StopChan: stopChan, JobChannel: &jobChannel, Action: action}
		w.spawn()
		(*p).Workers = append((*p).Workers, &w)
		log.Println("Started worker", i)
	}
}

// Dispatch will assign a job to one of the workers
func (p *Pool) Dispatch(job Job) {
	pool := *p
	*(pool.JobChannel) <- job
}

// Stop will kill all the workers in the pool
func (p *Pool) Stop() {
	pool := *p
	for i := range pool.Workers {
		w := *(pool.Workers[i])
		w.StopChan <- true
		log.Println("Stopped worker", i)
	}
}
