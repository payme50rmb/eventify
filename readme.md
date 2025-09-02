# Eventify

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Eventify is a lightweight, thread-safe event emitter for Go that simplifies building event-driven applications with support for both synchronous and asynchronous event handling.

## Features

- Simple and intuitive API
- Support for both synchronous and asynchronous event handling
- Thread-safe operations
- Named listeners for selective unregistering
- Built-in error handling
- Extensible event and listener interfaces

## Installation

```bash
go get github.com/payme50rmb/eventify
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/payme50rmb/eventify"
)

func main() {
    // Create a new Eventify instance with no logging
    ev := eventify.NewEventify(&eventify.NoLog{})
    
    // Register a simple event listener
    ev.Register("user.created", eventify.NewListener(func(event eventify.Event) {
        fmt.Printf("New user created: %s\n", string(event.Payload()))
    }))
    
    // Emit an event
    ev.Emit("user.created", []byte("john.doe@example.com"))
}

```

## Core Concepts

### Events

Events in Eventify implement the `Event` interface. You can use the built-in event type or create custom events.

### Listeners

Listeners handle events and can be either function-based or struct-based. They can be named for selective unregistering.

## Basic Usage

### Registering and Emitting Events

```go
package main

import (
    "fmt"
    "time"
    "github.com/payme50rmb/eventify"
)

func main() {
    // Create a new Eventify instance
    ev := eventify.NewEventify(&eventify.NoLog{})
    
    // Channel to collect event data
    dataCh := make(chan []byte, 10) // Buffered channel to prevent blocking
    
    // Register a listener for "user.login" events
    loginListener := eventify.NewNamedListener("login-tracker", func(event eventify.Event) {
        dataCh <- []byte(fmt.Sprintf("Login event: %s", string(event.Payload())))
    })
    ev.Register("user.login", loginListener)
    
    // Register another listener for the same event
    ev.Register("user.login", eventify.NewListener(func(event eventify.Event) {
        dataCh <- []byte(fmt.Sprintf("Another login handler: %s", string(event.Payload())))
    }))
    
    // Unregister the named listener after 5 seconds
    go func() {
        time.Sleep(5 * time.Second)
        ev.Unregister("user.login", loginListener)
        dataCh <- []byte("Login tracker unregistered")
    }()
    
    // Emit events in a separate goroutine
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        
        for i := 0; i < 10; i++ {
            select {
            case <-ticker.C:
                ev.Emit("user.login", []byte(fmt.Sprintf("User %d logged in", i)))
            }
        }
        close(dataCh) // Close channel when done
    }()
    
    // Process events from the channel
    for data := range dataCh {
        fmt.Println(string(data))
    }
}

```

## Advanced Usage

### Creating Custom Events

You can create custom events by implementing the `eventify.Event` interface. Here's an example of a custom event with error handling:

```go
// MyEvent implements the eventify.Event interface
type MyEvent struct {
    // eventify.IAmAsync // Uncomment to make this event async
    // Add any custom fields here
}

func (e *MyEvent) Type() string {
    return "my_event"
}

func (e *MyEvent) Payload() []byte {
    return []byte("custom payload")
}

// Optional: Implement ErrorHandler to handle errors from listeners
func (e *MyEvent) ErrorHandler(event eventify.Event, err error) {
    log.Printf("Error handling event %s: %v", event.Type(), err)
}
```

### Creating Custom Listeners

Create custom listeners by implementing the `Listener` interface:

```go
// MyListener implements the Listener interface
type MyListener struct {
    Name string
    // eventify.IAmAsync // Uncomment to make this listener async
}

func (l *MyListener) Name() string {
    return l.Name
}

func (l *MyListener) Handle(event eventify.Event) error {
    // Handle the event
    fmt.Printf("Handling event: %s with payload: %s\n", 
        event.Type(), 
        string(event.Payload()))
    return nil // Return error if handling fails
}
```

## Examples

### Registering and Unregistering Listeners

```go
package main

import (
	"fmt"
	"errors"

	"github.com/payme50rmb/eventify"
)
// MyEvent is the event of the eventify
type MyEvent struct {
	// eventify.IAmAsync # This event will run as async
}
// Type is the type of the event
func (e *MyEvent) Type() string {
	return "my_event"
}

// Payload is the payload of the event
func (e *MyEvent) Payload() []byte {
	return []byte("my_event")
}

// ErrorHandler is called when an error occurs in the listener
// The ErrorHandler is optional.
// If the ErrorHandler is not set, the error will be ignored.
func (e *MyEvent) ErrorHandler(event eventify.Event, err error) {
	fmt.Println(event)
	fmt.Println(err)
}

// MyListener is the listener of the event
type MyListener struct {
	// eventify.IAmAsync # This listener will run as async
}

// Name is the name of the listener, it can be unregistered by the name
// If do not set the name, the listener can not be unregistered.
func (l *MyListener) Name() string {
	return "my_listener"
}

// Handle is the handler of the listener
func (l *MyListener) Handle(event eventify.Event) error {
	fmt.Println(event)
	return errors.New("error")
}

```
