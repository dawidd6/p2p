// Package worker implements a worker pool
package worker

import (
	"sync"
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
			defer pool.waitGroup.Done()

			for job := range pool.channel {
				job()
			}
		}()
	}
}

// Enqueue sends a job to free worker
func (pool *Pool) Enqueue(job func()) {
	pool.channel <- job
}

// Finish closes the channel and waits for all workers to end
func (pool *Pool) Finish() {
	if pool.channelOpen {
		close(pool.channel)
		pool.channelOpen = false
		pool.waitGroup.Wait()
	}
}
