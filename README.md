# Generic PubSub Implementation in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/mdigger/pubsub.svg)](https://pkg.go.dev/github.com/mdigger/pubsub)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A thread-safe, generic Publish-Subscribe implementation in Go with blocking semantics.

## Features

- **Type-safe generics** - Works with any comparable key type and any message type
- **Thread-safe** - Safe for concurrent use by multiple goroutines
- **Blocking semantics** - Guaranteed message delivery (when channels are properly managed)
- **Lightweight** - Minimal dependencies (only standard library)
- **Efficient** - O(1) subscription lookups and O(n) publishes (n = subscribers per key)

## Installation

```bash
go get github.com/mdigger/pubsub
```

## Usage

```go
import "github.com/mdigger/pubsub"

// Create a new PubSub instance
ps := pubsub.New[string, string]()

// Create subscriber channels
ch1 := make(chan string, 10)
ch2 := make(chan string, 10)

// Subscribe channels to topics
ps.Subscribe([]string{"topic1", "topic2"}, ch1)
ps.Subscribe([]string{"topic1"}, ch2)

// Publish messages (blocks until all subscribers receive)
go func() {
    ps.Publish("topic1", "hello world")
}()

// Receive messages
msg := <-ch1
fmt.Println("Received:", msg)

// Unsubscribe when done
ps.Unsubscribe([]string{"topic1"}, ch1)
```

## Performance Considerations

1. **Channel Buffering**: Use buffered channels to prevent blocking publishers
2. **Key Cardinality**: Many unique keys will increase memory usage
3. **Fan-out**: Publishing to keys with many subscribers will be slower

## Best Practices

1. Always use buffered channels with sufficient capacity
2. Ensure subscribers are actively reading from channels
3. Consider using separate PubSub instances for different domains
4. Clean up unused subscriptions with Unsubscribe

## Alternatives

For non-blocking semantics or different delivery guarantees, consider:
- [go-redis PubSub](https://redis.io/topics/pubsub)
- [NATS](https://nats.io/)
- [Sarama](https://github.com/Shopify/sarama) (Kafka)
