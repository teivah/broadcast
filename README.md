# broadcast

![CI](https://github.com/teivah/broadcast/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/teivah/broadcast)](https://goreportcard.com/report/github.com/teivah/broadcast)

A broadcasting library for Go.

## What?

`broadcast` is a library that allows sending repeated notifications to multiple goroutines with guaranteed delivery.

## Why?

### Why not Channels?

The standard way to handle notifications is via a `chan struct{}`. However, sending a message to a channel is received by a single goroutine. 

The only operation that is broadcasted to multiple goroutines is a channel closure. Yet, if the channel is closed, there's no way to send a message again.

❌ Repeated notifications to multiple goroutines
✅ Guaranteed delivery

### Why not sync.Cond?

`sync.Cond` is the standard solution based on condition variables to set up containers of goroutines waiting for a specific condition.

There's one caveat to keep in mind, though: the `Broadcast()` method doesn't guarantee that a goroutine will receive the notification. Indeed, the notification will be lost if the listener goroutine isn't waiting on the `Wait()` method.

✅ Repeated notifications to multiple goroutines
❌ Guaranteed delivery

## How?

### Step by Step

First, we need to create a `Relay`:

```go
relay := broadcast.NewRelay() // Create a new relay
```

Once a `Relay` is created, we can create a new listener using the `Listen` method. As the `broadcast` library relies internally on channels, it accepts a capacity:

````go
list := relay.Listener(1) // Create a new listener based on a channel with a one capacity
````

A `Relay` can send a notification in two different manners:
* Blocking: using the `Notify` method
* Non-blocking: using the `Broadcast` method

The main difference is that `Notify` guarantees that a listener will receive the notification. This method blocks until the notification have been sent to all the `Listener`. Conversely, `Broadcast` isn't blocking:

```go
relay.Notify() // Send a blocking notification (delivery is guaranteed)
relay.Broadcast() // Send a non-blocking notification (delivery is not guaranteed)
```

On the `Listener` side, we can access the internal channel using `Ch`:

```go
<-list.Ch() // Wait on a notification
```

We can close a `Listener` using `Close`:

```go
list.Close() // Close a listener
```

Please note that this operation is purposely not thread-safe. It can be called multiple times by the same goroutine but shouldn't be called in parallel by multiple goroutines as it can lead to a panic (trying to close the same channel multiple times).

Last but not least, we can close a `Relay` using `Close`:

```go
relay.Close() // Close a relay
```

### Example

```go
relay := broadcast.NewRelay() // Create a relay
defer relay.Close()

list1 := relay.Listener(1) // Create a listener with one capacity
defer list1.Close()
list2 := relay.Listener(1) // Create a listener with one capacity
defer list2.Close()

// Listener goroutines
f := func(i int, list *broadcast.Listener) {
	for range list.Ch() { // Waits for receiving notifications
		fmt.Printf("listener %d has received a notification\n", i)
	}
}
go f(1, list1)
go f(2, list2)

// Notifier goroutine
for i := 0; i < 5; i++ {
	time.Sleep(time.Second)
	relay.Notify() // Send notifications with guaranteed delivery
}
for i := 0; i < 5; i++ {
	time.Sleep(time.Second)
	relay.Broadcast() // Send notifications without guaranteed delivery
}
```
