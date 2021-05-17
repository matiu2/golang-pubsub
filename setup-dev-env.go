package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

/// Creates the pubub topic if it doesn't exist
/// Panics on fail so no error is returned
func EnsureTopic(client *pubsub.Client, projectID, topicID string) {
	ctx := context.Background()

	// See if our topic already exists
	topic := client.Topic(topicID)

	// Create the topic if needed
	exists, err := topic.Exists(ctx)
	if !exists || err != nil {
		t, err := client.CreateTopic(ctx, topicID)
		if err != nil {
			// Probably it already exists, so just log it and return
			log.Printf("Error creating topic: %v", err)
		}
		topic = t
		log.Printf("Topic created: %v\n", topic)
	} else {
		log.Printf("Topic found: %v\n", topic)
	}
}

func EnsureSubscription(client *pubsub.Client, topicID, subID string) {
	ctx := context.Background()

	// List the subscriptions in case it's already there
	sub := client.Subscription(subID)
	exists, err := sub.Exists(ctx)
	if !exists || err != nil {
		// Create the subscription
		s, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:       client.Topic(topicID),
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			// Probably it already exists, so just log it and return
			log.Printf("Error creating subscription: %v", err)
		}
		sub = s
		log.Printf("Subscription created: %v", sub)
	} else {
		log.Printf("Subscrition found: %v", sub)
	}
}
