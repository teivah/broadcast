# broadcast

![CI](https://github.com/teivah/broadcast/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/teivah/broadcast)](https://goreportcard.com/report/github.com/teivah/broadcast)

A broadcasting library for Go.

## What?

`broadcast` is a library that allows sending repeated notifications to multiple goroutines with possible guaranteed delivery.

## Why?

### Why Not Using Channels?

Sending a message to a channel is received by a single goroutine. The only operation that is broadcasted to multiple goroutines is a channel closure.

However, if the channel is closed, there's no way to send a message again.

### Why Not Using sync.Cond?

`sync.Cond` is the standard solution based on condition variables to set up containers of goroutines waiting for a specific condition. There's one caveat to keep in mind, though: the `Broadcast()` method doesn't guarantee that a goroutine will receive the notification. Indeed, the notification will be lost if the listener goroutine isn't waiting on the `Wait()` method.

## How?

First, we need to create a `Relay`:

```go
relay, err := broadcast.NewRelay() // Create a new relay
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

Last but not least, we can close a `Listener` using `Close`:

```go
list.Close()
```

Please note that this operation is purposely not thread-safe. It can be called multiple times by the same goroutine but shouldn't be called in parallel by multiple goroutines as it can lead to a panic (trying to close the same channel multiple times).
