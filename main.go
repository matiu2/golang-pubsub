package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
)

func pullMsgs(projectID, subID string, done *sync.WaitGroup) {
	defer done.Done()
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	log.Printf("Listening for messages: %s - %s", projectID, subID)

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		log.Printf("Got message: %q\n", string(msg.Data))
		msg.Ack()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Receive: %v", err)
	}
}

func publishMsgs(projectID string, topicID string, n int, done *sync.WaitGroup) {
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

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
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
		log.Fatalf("%d of %d messages did not publish successfully", totalErrors, n)
	}
}

func main() {
	godotenv.Load(".env")
	projectID := "test"
	topicID := "test-topic"
	subID := "test-sub"

	// Create the pubsub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	// Make the pubsub topic (if needed)
	EnsureTopic(client, projectID, topicID)
	EnsureSubscription(client, topicID, subID)

	var wg sync.WaitGroup
	wg.Add(2)

	go pullMsgs(projectID, subID, &wg)
	go publishMsgs(projectID, topicID, 10, &wg)

	log.Printf("Waiting for all the messages")
	wg.Wait()

}
