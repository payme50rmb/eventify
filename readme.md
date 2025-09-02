# Eventify

[![Go Reference](https://pkg.go.dev/badge/github.com/payme50rmb/eventify.svg)](https://pkg.go.dev/github.com/payme50rmb/eventify)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/payme50rmb/eventify)](https://goreportcard.com/report/github.com/payme50rmb/eventify)

Eventify is a lightweight, thread-safe event emitter for Go that simplifies building event-driven applications with support for both synchronous and asynchronous event handling.

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🚀 Simple API | Intuitive and easy-to-use interface |
| ⚡ Sync/Async | Support for both synchronous and asynchronous event handling |
| 🔒 Thread-safe | Safe for concurrent use |
| 🏷️ Named Listeners | Register and unregister listeners by name |
| 🛡️ Error Handling | Built-in error handling capabilities |
| 🔌 Extensible | Create custom events and listeners |

## 🚀 Installation

```bash
go get github.com/payme50rmb/eventify
```

## ⚡ Quick Start

```go
package main

import (
	"fmt"
	"github.com/payme50rmb/eventify"
)

func main() {
	// Initialize a new Eventify instance
	ev := eventify.New()

	// Register an event listener
	ev.Register("user.created", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("New user created: %s\n", string(event.Payload()))
		return nil
	}))

	// Emit an event
	ev.EmitBy("user.created", []byte("john.doe@example.com"))
}
```

## 🧩 Core Concepts

### Events

Events in Eventify implement the `Event` interface. You can use the built-in event type or create custom events by implementing:

```go
type Event interface {
    Type() string
    Payload() []byte
}
```

### Listeners

Listeners handle events and can be either function-based or struct-based. They implement:

```go
type Listener interface {
    Handle(event Event) error
}
```

## 🛠️ Basic Usage

### Registering and Emitting Events

```go
package main

import (
	"fmt"
	"github.com/payme50rmb/eventify"
	"time"
)

func main() {
	// Initialize Eventify
	ev := eventify.New()

	// Channel to collect event data
	dataCh := make(chan string, 10)

	// Create a named listener
	loginTracker := eventify.NewNamedListener("tracker", func(event eventify.Event) error {
		dataCh <- fmt.Sprintf("📊 Tracking login: %s", string(event.Payload()))
		return nil
	})

	// Register listeners
	ev.Register("user.login", loginTracker)
	ev.Register("user.login", eventify.NewListener(func(event eventify.Event) error {
		dataCh <- fmt.Sprintf("🔔 New login: %s", string(event.Payload()))
		return nil
	}))

	// Unregister after delay
	go func() {
		time.Sleep(3 * time.Second)
		ev.Unregister("user.login", loginTracker)
		dataCh <- "❌ Login tracker unregistered"
	}()

	// Emit events
	go func() {
		for i := 1; i <= 5; i++ {
			ev.EmitBy("user.login", []byte(fmt.Sprintf("user%d@example.com", i)))
			time.Sleep(time.Second)
		}
		close(dataCh)
	}()

	// Process events
	for msg := range dataCh {
		fmt.Println(msg)
	}
}

```

Example Output:

```
📊 Tracking login: user1@example.com
🔔 New login: user1@example.com
📊 Tracking login: user2@example.com
🔔 New login: user2@example.com
📊 Tracking login: user3@example.com
🔔 New login: user3@example.com
❌ Login tracker unregistered
🔔 New login: user4@example.com
🔔 New login: user5@example.com
```

## 📚 Examples

### 1. Basic Event Handling

This example shows how to set up basic event handling with Eventify:

```go
package main

import (
	"fmt"
	"github.com/payme50rmb/eventify"
)

func main() {
	// Initialize with default config
	ev := eventify.New()

	// Register event listeners
	ev.Register("user.created", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("👤 New user created: %s\n", string(event.Payload()))
		return nil
	}))

	ev.Register("order.placed", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("🛒 New order: %s\n", string(event.Payload()))
		return nil
	}))

	// Emit events
	ev.EmitBy("user.created", []byte("alice@example.com"))
	ev.EmitBy("order.placed", []byte("order_12345"))
}
```

### 2. Wildcard Event Matching

```go
package main

import (
	"fmt"
	"github.com/payme50rmb/eventify"
)

func main() {
	ev := eventify.New()

	// This will match any event starting with 'user.'
	ev.Register("user.*", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("🔔 User event (%s): %s\n", 
			event.Type(), 
			string(event.Payload()))
		return nil
	}))

	// These events will all trigger the above listener
	ev.EmitBy("user.created", []byte("bob@example.com"))
	ev.EmitBy("user.updated", []byte("profile updated"))
	ev.EmitBy("user.deleted", []byte("user_456"))

	// This will match any event ending with '.order.paid'
	ev.Register("*.order.paid", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("🔔 User event (%s): %s\n", 
			event.Type(), 
			string(event.Payload()))
		return nil
	}))

	// These events will all trigger the above listener
	ev.EmitBy("alipay.order.paid", []byte("order_12345"))
	ev.EmitBy("wechat.order.paid", []byte("order_12345"))

	// This will match any event containing 'payment'
	ev.Register("*.payment.*", eventify.NewListener(func(event eventify.Event) error {
		fmt.Printf("🔔 User event (%s): %s\n", 
			event.Type(), 
			string(event.Payload()))
		return nil
	}))

	// These events will all trigger the above listener
	ev.EmitBy("alipay.payment.success", []byte("order_12345"))
	ev.EmitBy("wechat.payment.success", []byte("order_12345"))
	ev.EmitBy("alipay.payment.failed", []byte("order_12345"))
	ev.EmitBy("wechat.payment.failed", []byte("order_12345"))
}

```

### 3. Error Handling

```go
package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/payme50rmb/eventify"
)

type PaymentEvent struct {
	Amount int64
	Error  chan error
}

func (e *PaymentEvent) Type() string { return "payment.processed" }
func (e *PaymentEvent) Payload() []byte {
	v := big.NewInt(e.Amount)
	return v.Bytes()
}

func (e *PaymentEvent) ErrorHandler(_ eventify.Event, err error) {
	// Add your error handling logic here
	// e.g., retry, notify, or update status
	e.Error <- err
	close(e.Error)
}

func main() {
	ev := eventify.New()

	ev.Register("payment.processed", eventify.NewListener(func(event eventify.Event) error {
		amount := big.NewInt(0).SetBytes(event.Payload())
		fa := amount.Int64()
		if fa > 90 {
			return errors.New("payment amount too high")
		}
		return nil
	}))
	errCh := make(chan error, 1)
	// This will trigger the error handler
	ev.Emit(&PaymentEvent{Amount: 100, Error: errCh})

	for err := range errCh {
		fmt.Println("❌", err.Error())
	}
}

```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Report bugs** by opening an issue
2. **Suggest features** or improvements
3. **Submit pull requests**

### Development Setup

1. Fork the repository
2. Clone your fork: 
   ```bash
   git clone https://github.com/your-username/eventify.git
   cd eventify
   ```
3. Run tests:
   ```bash
   go test -v ./...
   ```
4. Make your changes and ensure tests pass
5. Submit a pull request

### Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `gofmt` and `golint` before submitting PRs
- Write tests for new features and bug fixes

## 🌟 Show Your Support

Give a ⭐️ if this project helped you!

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
