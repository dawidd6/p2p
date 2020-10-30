package worker_test

import (
	"sync/atomic"
	"testing"

	"github.com/dawidd6/p2p/pkg/worker"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	got := uint32(0)
	expected := uint32(10)
	pool := worker.New(2)
	pool.Start()
	for i := uint32(0); i < expected; i++ {
		pool.Enqueue(func() {
			atomic.AddUint32(&got, 1)
		})
	}
	pool.Stop()

	assert.Equal(t, expected, got)
}
