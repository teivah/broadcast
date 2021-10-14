// Package broadcast allows to send repeated notifications to multiple goroutines.
package broadcast

import (
	"sync"
)

// Relay is the struct in charge of handling the listeners and dispatching the notifications.
type Relay struct {
	mu      sync.RWMutex
	n       uint32
	clients map[uint32]*Listener
}

// NewRelay is the factory to create a Relay.
func NewRelay() (*Relay, error) {
	return &Relay{
		clients: make(map[uint32]*Listener),
	}, nil
}

// Notify sends a notification to all the listeners.
// It guarantees that all the listeners will receive the notification.
func (r *Relay) Notify() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		client.ch <- struct{}{}
	}
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
		ch:        make(chan struct{}, capacity),
		id:        r.n,
		broadcast: r,
	}
	r.clients[r.n] = listener
	r.n++
	return listener
}

func (r *Relay) close(n uint32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, n)
}

// Listener is a Relay listener.
type Listener struct {
	ch        chan struct{}
	id        uint32
	broadcast *Relay
	closed    bool
}

// Ch returns the Listener channel.
func (l *Listener) Ch() <-chan struct{} {
	return l.ch
}

// Close closes a listener.
// This operation isn't thread safe.
func (l *Listener) Close() {
	if l.closed {
		return
	}

	l.broadcast.close(l.id)
	close(l.ch)
	l.closed = true
}
