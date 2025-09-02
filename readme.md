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

### 1. Basic Event Handling

```go
package main

import (
	"fmt"
	"github.com/payme50rmb/eventify"
)

func main() {
	ev := eventify.NewEventify(nil)

	// Register a simple event listener
	ev.Register("user.created", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("New user created: %s\n", string(event.Payload()))
		return nil
	}))

	// Emit an event
	ev.Emit("user.created", []byte("john@example.com"))
}
```

### 2. Error Handling

```go
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/payme50rmb/eventify"
)

type PaymentEvent struct {
	Amount int
}

func (e *PaymentEvent) Type() string    { return "payment.processed" }
func (e *PaymentEvent) Payload() []byte { return []byte(fmt.Sprintf(`{"amount":%d}`, e.Amount)) }
func (e *PaymentEvent) ErrorHandler(_ eventify.Event, err error) {
	log.Printf("Error processing payment: %v", err)
}

func main() {
	ev := eventify.NewEventify(nil)

	ev.Register("payment.processed", eventify.NewListener(func(event eventify.Event) error {
		amount := 100 // Extract from event.Payload() in real code
		if amount > 90 {
			return errors.New("payment amount too high")
		}
		return nil
	}))

	// This will trigger the error handler
	ev._Emit(&PaymentEvent{Amount: 100})
}
```

### 3. Middleware Pattern

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/payme50rmb/eventify"
)

func loggingMiddleware(next eventify.Listener) eventify.Listener {
	return eventify.NewListener(func(event eventify.Event) error {
		start := time.Now()
		defer func() {
			log.Printf("Event %s processed in %v", event.Type(), time.Since(start))
		}()
		return next.Handle(event)
	})
}

func main() {
	ev := eventify.NewEventify(nil)

	eventHandler := eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("Processing event: %s\n", event.Type())
		time.Sleep(100 * time.Millisecond) // Simulate work
		return nil
	})

	// Wrap the handler with middleware
	ev.Register("order.placed", loggingMiddleware(eventHandler))

	ev.Emit("order.placed", []byte("order data"))
}
```

### 4. Request-Reply Pattern

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/payme50rmb/eventify"
)

type Request struct {
	ID     string
	Method string
	Data   interface{}
}

type Response struct {
	RequestID string
	Result    interface{}
	Error     string
}

func main() {
	ev := eventify.NewEventify(nil)

	// Request handler
	ev.Register("api.getUser", eventify.NewListener(func(event eventify.Event) error {
		var req Request
		if err := json.Unmarshal(event.Payload(), &req); err != nil {
			return err
		}

		// Process request and prepare response
		response := Response{
			RequestID: req.ID,
			Result:    map[string]string{"id": "123", "name": "John Doe"},
		}

		// Send response back
		respData, _ := json.Marshal(response)
		ev.Emit("api.response."+req.ID, respData)
		return nil
	}))

	// Make a request
	reqID := "req_123"
	req := Request{
		ID:     reqID,
		Method: "getUser",
		Data:   map[string]string{"userID": "123"},
	}
	reqData, _ := json.Marshal(req)

	// Set up response handler
	respCh := make(chan []byte, 1)
	ev.Register("api.response."+reqID, eventify.NewListener(func(event eventify.Event) error {
		respCh <- event.Payload()
		return nil
	}))

	// Send request
	ev.Emit("api.getUser", reqData)

	// Wait for response
	respData := <-respCh
	var resp Response
	json.Unmarshal(respData, &resp)
	fmt.Printf("Got response: %+v\n", resp)
}
```

### 5. Event Batching

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/payme50rmb/eventify"
)

type BatchProcessor struct {
	ev      *eventify.Eventify
	buffer  []eventify.Event
	mu      sync.Mutex
	ticker  *time.Ticker
	maxSize int
}

func NewBatchProcessor(ev *eventify.Eventify, interval time.Duration, maxSize int) *BatchProcessor {
	bp := &BatchProcessor{
		ev:      ev,
		buffer:  make([]eventify.Event, 0, maxSize),
		ticker:  time.NewTicker(interval),
		maxSize: maxSize,
	}
	go bp.start()
	return bp
}

func (bp *BatchProcessor) start() {
	for range bp.ticker.C {
		bp.processBatch()
	}
}

func (bp *BatchProcessor) processBatch() {
	bp.mu.Lock()
	if len(bp.buffer) == 0 {
		bp.mu.Unlock()
		return
	}

	batch := make([]eventify.Event, len(bp.buffer))
	copy(batch, bp.buffer)
	bp.buffer = bp.buffer[:0]
	bp.mu.Unlock()

	// Process the batch
	fmt.Printf("Processing batch of %d events\n", len(batch))
	for _, event := range batch {
		fmt.Printf(" - %s: %s\n", event.Type(), string(event.Payload()))
	}
}

func (bp *BatchProcessor) AddEvent(event eventify.Event) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.buffer = append(bp.buffer, event)
	if len(bp.buffer) >= bp.maxSize {
		bp.processBatch()
	}
}

func main() {
	ev := eventify.NewEventify(nil)
	bp := NewBatchProcessor(ev, 5*time.Second, 10)

	// Register event handler that adds to batch
	ev.Register("log.event", eventify.NewListener(func(event eventify.Event) error {
		bp.AddEvent(event)
		return nil
	}))

	// Generate some events
	for i := 0; i < 15; i++ {
		ev.Emit("log.event", []byte(fmt.Sprintf("Event %d", i)))
		time.Sleep(500 * time.Millisecond)
	}

	// Wait for any remaining events to be processed
	time.Sleep(6 * time.Second)
}
```

### 6. Graceful Shutdown

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/payme50rmb/eventify"
)

func main() {
	ev := eventify.NewEventify(nil)

	// Simulate a long-running event processor
	ev.Register("data.process", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("Processing: %s\n", string(event.Payload()))
		time.Sleep(1 * time.Second) // Simulate work
		return nil
	}))

	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start emitting events in the background
	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				ev.Emit("data.process", []byte(fmt.Sprintf("Data %d", i)))
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	fmt.Println("\nShutting down gracefully...")

	// Give some time for in-flight events to complete
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("Shutdown complete")
	}
}
```

### 7. Registering and Unregistering Listeners

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
