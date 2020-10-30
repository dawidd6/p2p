package notifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dawidd6/p2p/pkg/notifier"
)

func TestNotBlocking(t *testing.T) {
	notify := notifier.NewNotBlocking()
	start := make(chan struct{})
	stop := make(chan struct{})
	count := 0

	go func() {
		start <- struct{}{}
		count = 1
		stop <- struct{}{}
		notify.Notify()
		notify.Notify()
		notify.Notify()
		start <- struct{}{}
		count = 2
		stop <- struct{}{}
	}()

	<-start
	<-stop
	assert.Equal(t, 1, count)

	<-start
	<-stop
	assert.Equal(t, 2, count)

	<-notify.Wait()
}

func TestBlocking(t *testing.T) {
	notify := notifier.NewBlocking()
	start := make(chan struct{})
	stop := make(chan struct{})
	count := 0

	go func() {
		start <- struct{}{}
		count = 1
		stop <- struct{}{}
		notify.Notify()
		start <- struct{}{}
		count = 2
		stop <- struct{}{}
	}()

	<-start
	<-stop
	assert.Equal(t, 1, count)

	<-notify.Wait()

	<-start
	<-stop
	assert.Equal(t, 2, count)
}
