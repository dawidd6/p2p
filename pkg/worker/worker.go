package worker

import (
	"sync"
)

type Pool struct {
	waitGroup   sync.WaitGroup
	channel     chan func()
	channelOpen bool
	count       int
}

func New(count int) *Pool {
	return &Pool{
		count: count,
	}
}

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

func (pool *Pool) Enqueue(job func()) {
	pool.channel <- job
}

func (pool *Pool) Finish() {
	if pool.channelOpen {
		close(pool.channel)
		pool.channelOpen = false
		pool.waitGroup.Wait()
	}
}
