// Package broadcast allows to send repeated notifications to multiple goroutines.
package broadcast

import (
	"context"
	"sync"
)

// Relay is the struct in charge of handling the listeners and dispatching the notifications.
type Relay struct {
	mu      sync.RWMutex
	n       uint32
	clients map[uint32]*Listener
}

// NewRelay is the factory to create a Relay.
func NewRelay() *Relay {
	return &Relay{
		clients: make(map[uint32]*Listener),
	}
}

// Notify sends a notification to all the listeners.
// It guarantees that all the listeners will receive the notification.
func (r *Relay) Notify() {
	_ = r.NotifyCtx(context.Background()) //lint:ignore
}

// NotifyCtx tries sending a notification to all the listeners until the context times out or is canceled.
func (r *Relay) NotifyCtx(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		select {
		case client.ch <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// Broadcast broadcasts a notification to all the listeners.
// The notification is sent in a non-blocking manner, so there's no guarantee that a listener receives it.
func (r *Relay) Broadcast() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		select {
		case client.ch <- struct{}{}:
		default:
		}
	}
}

// Listener creates a new listener given a channel capacity.
func (r *Relay) Listener(capacity int) *Listener {
	r.mu.Lock()
	defer r.mu.Unlock()

	listener := &Listener{
		ch:    make(chan struct{}, capacity),
		id:    r.n,
		relay: r,
	}
	r.clients[r.n] = listener
	r.n++
	return listener
}

// Close closes a relay.
// This operation can be safely called in the meantime as Listener.Close()
func (r *Relay) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, client := range r.clients {
		client.Close()
	}
	r.clients = make(map[uint32]*Listener)
}

func (r *Relay) removeListener(l *Listener) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clients, l.id)
}

// Listener is a Relay listener.
type Listener struct {
	ch    chan struct{}
	id    uint32
	relay *Relay
	once  sync.Once
}

// Ch returns the Listener channel.
func (l *Listener) Ch() <-chan struct{} {
	return l.ch
}

// Close closes a listener.
// This operation can be safely called in the meantime as Relay.Close()
func (l *Listener) Close() {
	l.once.Do(func() {
		close(l.ch)
		l.relay.removeListener(l)
	})
}
