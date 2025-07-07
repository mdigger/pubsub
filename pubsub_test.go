package pubsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/mdigger/pubsub"
)

func TestPublishWithContext(t *testing.T) {
	ps := pubsub.New[string, string]()
	ch := make(chan string) // unbuffered
	ps.Subscribe([]string{"topic"}, ch)

	t.Run("successful publish", func(t *testing.T) {
		ctx := context.Background()
		go func() {
			time.Sleep(10 * time.Millisecond)
			<-ch // consume the message
		}()

		delivered, err := ps.Publish(ctx, "topic", "msg")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if delivered != 1 {
			t.Errorf("expected 1 delivery, got %d", delivered)
		}
	})

	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		delivered, err := ps.Publish(ctx, "topic", "msg")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
		if delivered != 0 {
			t.Errorf("expected 0 deliveries, got %d", delivered)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		delivered, err := ps.Publish(ctx, "topic", "msg")
		if err != context.DeadlineExceeded {
			t.Errorf("expected DeadlineExceeded, got %v", err)
		}
		if delivered != 0 {
			t.Errorf("expected 0 deliveries, got %d", delivered)
		}
	})
}

func TestPublishWithTimeout(t *testing.T) {
	ps := pubsub.New[string, string]()
	ch := make(chan string) // unbuffered
	ps.Subscribe([]string{"topic"}, ch)

	delivered, err := ps.PublishWithTimeout("topic", "msg", 10*time.Millisecond)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
	if delivered != 0 {
		t.Errorf("expected 0 deliveries, got %d", delivered)
	}
}
