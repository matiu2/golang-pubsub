package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
)

func createTopic(projectID, topicID string) error {
	ctx := context.Background()
	log.Printf("Creating topic: %s - %s", projectID, topicID)
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	t, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		log.Fatalf("CreateTopic: %v", err)
	}
	log.Printf("Topic created: %v\n", t)
	return nil
}

func main() {
	godotenv.Load(".env")
	projectID := "test"
	topicID := "test-topic"
	createTopic(projectID, topicID)
}
