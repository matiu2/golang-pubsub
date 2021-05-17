package main

import (
	"context"
	"log"
	"sync"

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
	go PublishMsgs(projectID, topicID, 10, &wg)

	log.Printf("Waiting for all the messages")
	wg.Wait()

}
