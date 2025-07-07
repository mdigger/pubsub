// Package pubsub implements a generic Publish-Subscribe pattern.
// It allows publishers to send messages to multiple subscribers
// based on topic keys, with thread-safe operations.
// The implementation is generic, supporting any comparable key type
// and any message type.
package pubsub

import (
	"sync"
)

// PubSub implements the Publish-Subscribe pattern.
// It maintains a mapping of keys to subscriber channels,
// allowing efficient message distribution.
// K is the key type (must be comparable), T is the message type.
type PubSub[K comparable, T any] struct {
	mu          sync.RWMutex // protects subscribers map
	subscribers map[K]map[chan T]struct{}
}

// New creates and returns a new PubSub instance.
// The returned PubSub is ready to use with zero values initialized.
func New[K comparable, T any]() *PubSub[K, T] {
	return &PubSub[K, T]{
		subscribers: make(map[K]map[chan T]struct{}),
	}
}

// Subscribe adds a channel to receive messages for the specified keys.
// The channel will receive all messages published to any of the provided keys.
// If the channel is already subscribed to a key, this is a no-op.
//
// Note: The channel should have sufficient buffer space or active readers
// to prevent indefinite blocking in the Publish method.
func (ps *PubSub[K, T]) Subscribe(keys []K, ch chan T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, key := range keys {
		if _, exists := ps.subscribers[key]; !exists {
			ps.subscribers[key] = make(map[chan T]struct{})
		}

		ps.subscribers[key][ch] = struct{}{}
	}
}

// Unsubscribe removes a channel from receiving messages for the specified keys.
// After this call, the channel will no longer receive messages for these keys.
// If the channel wasn't subscribed to a key, that key is skipped.
// If all channels are unsubscribed from a key, the key is removed from the registry.
func (ps *PubSub[K, T]) Unsubscribe(keys []K, ch chan T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, key := range keys {
		if subs, exists := ps.subscribers[key]; exists {
			delete(subs, ch)

			if len(subs) == 0 {
				delete(ps.subscribers, key)
			}
		}
	}
}

// Publish sends a message to all channels subscribed to the specified key.
// This is a blocking operation - it will wait for each subscriber's channel
// to accept the message. The message is delivered to all subscribers in an
// unspecified order.
//
// Warning: If any subscriber's channel is not being read from, this will
// block indefinitely. Ensure all subscribed channels have active readers
// or sufficient buffer space.
func (ps *PubSub[K, T]) Publish(key K, msg T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if subs, exists := ps.subscribers[key]; exists {
		for ch := range subs {
			ch <- msg // This will block until the message is accepted
		}
	}
}
