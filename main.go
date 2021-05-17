package main

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
)

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

	messages := GenerateTestMessages(10)

	go PullMsgs(projectID, subID, &wg)
	go PublishMsgs(projectID, topicID, messages, &wg)

	log.Printf("Waiting for all the messages")
	wg.Wait()

}
