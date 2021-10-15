package broadcast_test

import (
	"context"
	"sync"
	"testing"

	"github.com/teivah/broadcast"
)

func TestNotify(t *testing.T) {
	relay := broadcast.NewRelay()
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

	relay.Close()
}

func TestBroadcast(t *testing.T) {
	relay := broadcast.NewRelay()
	relay.Listener(0)

	relay.Broadcast()
	relay.Close()
}

func TestNotifyCtx(t *testing.T) {
	relay := broadcast.NewRelay()
	relay.Listener(0)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	relay.NotifyCtx(ctx)
}

func TestRaceBroadcast(t *testing.T) {
	relay := broadcast.NewRelay()
	wg := sync.WaitGroup{}

	funcs := []func(){
		func() {
			listener := relay.Listener(1)
			select {
			case <-listener.Ch():
			default:
			}
			listener.Close()
		},
		func() {
			relay.Broadcast()
		},
	}

	for i := 0; i < 1000; i++ {
		for _, f := range funcs {
			f := f
			wg.Add(1)
			go func() {
				f()
				wg.Done()
			}()
		}
	}

	relay.Broadcast()
	wg.Wait()
	relay.Close()
}

func TestRaceClosure(t *testing.T) {
	for i := 0; i < 1000; i++ {
		relay := broadcast.NewRelay()
		list := []*broadcast.Listener{
			relay.Listener(0),
			relay.Listener(0),
			relay.Listener(0),
		}
		go func() {
			relay.Close()
		}()
		go func() {
			for _, l := range list {
				l.Close()
			}
		}()
	}
}
