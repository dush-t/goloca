package pool

// Job contains the data that a WorkerAction will consume
type Job struct {
	DataPoint    interface{}
	CompleteChan chan bool
}

// WorkerAction is a function that a worker can be registered to perform
type WorkerAction func(Job) error
