// Package worker implements a worker pool
package worker

import (
	"sync"
    "time"
)

// Pool holds variables needed to coordinate the workers
type Pool struct {
	waitGroup   sync.WaitGroup
	channel     chan func()
	channelOpen bool
	count       int
}

// New returns a new Pool instance
func New(count int) *Pool {
	return &Pool{
		count: count,
	}
}

// Start fires up workers waiting for jobs to arrive
func (pool *Pool) Start() {
	pool.channel = make(chan func(), pool.count)
	pool.channelOpen = true

	for i := 0; i < pool.count; i++ {
		pool.waitGroup.Add(1)
		go func() {
			for job := range pool.channel {
                time.Sleep(time.Millisecond*100)
				job()
			}
			pool.waitGroup.Done()
		}()
	}
}

// Enqueue sends a job to free worker
func (pool *Pool) Enqueue(job func()) {
	pool.channel <- job
}

// Stop closes the channel and waits for all workers to end
func (pool *Pool) Stop() {
	if pool.channelOpen {
		close(pool.channel)
		pool.channelOpen = false
		pool.waitGroup.Wait()
	}
}
