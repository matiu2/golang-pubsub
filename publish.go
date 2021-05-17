/// pubsub message publisher - used for integration testing
package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

func PublishMsgs(projectID string, topicID string, messages []Message, done *sync.WaitGroup) {
	defer done.Done()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)

	for i, msg := range messages {
		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("Unable to convert outgoing message into json: %v", msg)
		}
		result := t.Publish(ctx, &pubsub.Message{
			Data: bytes,
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				log.Printf("Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			log.Printf("Published message %d; msg ID: %v\n", i, id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		log.Fatalf("%d of %d messages did not publish successfully", totalErrors, len(messages))
	}
}
