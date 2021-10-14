package broadcast_test

import (
	"broadcast"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotify(t *testing.T) {
	relay, err := broadcast.NewRelay()
	require.NoError(t, err)
	list1 := relay.Listener(0)
	list2 := relay.Listener(0)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		<-list1.Ch()
		wg.Done()
	}()
	go func() {
		<-list2.Ch()
		wg.Done()
	}()
	relay.Notify()
	wg.Wait()

	wg.Add(1)
	go func() {
		<-list2.Ch()
		wg.Done()
	}()
	list1.Close()
	relay.Notify()
	wg.Wait()

	list2.Close()
	relay.Notify()

	list1.Close()
	list2.Close()
}

func TestBroadcast(t *testing.T) {
	relay, err := broadcast.NewRelay()
	require.NoError(t, err)
	relay.Listener(0)

	relay.Broadcast()
}
