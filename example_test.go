package pubsub_test

import (
	"context"
	"fmt"
	"time"

	"github.com/mdigger/pubsub"
)

func Example_contextAwarePublishing() {
	ps := pubsub.New[string, string]()

	// Create a buffered channel to ensure the example works
	ch := make(chan string, 2)
	ps.Subscribe([]string{"alerts"}, ch)

	// Regular publish with background context
	delivered, err := ps.Publish(context.Background(), "alerts", "urgent")
	fmt.Printf("Delivered: %d, Error: %v\n", delivered, err)

	// Publish with timeout (using buffered channel to ensure delivery)
	delivered, err = ps.PublishWithTimeout("alerts", "timeout-test", 100*time.Millisecond)
	fmt.Printf("With timeout - Delivered: %d, Error: %v\n", delivered, err)

	// Consume the messages to clean up
	<-ch
	<-ch

	// Output:
	// Delivered: 1, Error: <nil>
	// With timeout - Delivered: 1, Error: <nil>
}
