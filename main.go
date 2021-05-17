package main

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

func pullMsgs(projectID, subID string) {
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

/// Creates the pubub topic if it doesn't exist
/// Panics on fail so no error is returned
func ensureTopic(client *pubsub.Client, projectID, topicID string) *pubsub.Topic {
	ctx := context.Background()

	// List the topics in case our topic already exists
	var topic *pubsub.Topic
	it := client.Topics(ctx)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("While listing topics: Next: %v", err)
		}
		log.Printf("Listing topic: %s", t.ID())
		if t.ID() == topicID {
			topic = t
			break
		}
	}

	// Create the topic if needed
	if topic == nil {
		t, err := client.CreateTopic(ctx, topicID)
		if err != nil {
			// Probably it already exists, so just log it and return
			log.Printf("Error creating topic: %v", err)
		}
		log.Printf("Topic created: %v\n", t)
		topic = t
	}

	return topic
}

func ensureSubscription(client *pubsub.Client, topic *pubsub.Topic, subID string) *pubsub.Subscription {
	ctx := context.Background()

	// List the subscriptions in case it's already there
	var sub *pubsub.Subscription
	it := client.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("While listing subscriptions: Next: %v", err)
		}
		log.Printf("Listing subscriptions: %s", s.ID())
		if s.ID() == subID {
			sub = s
			break
		}
	}

	if sub == nil {
		// Create the subscription
		s, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			// Probably it already exists, so just log it and return
			log.Printf("Error creating subscription: %v", err)
		}
		sub = s
		log.Printf("Subscription created: %v\n", sub)
	}

	return sub
}

func main() {
	godotenv.Load(".env")
	projectID := "test"
	topicID := "test-topic"
	subID := "test-sub"

	// Create the pubsub client
	ctx := context.Background()
	log.Printf("Creating topic: %s - %s", projectID, topicID)
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	// Make the pubsub topic (if needed)
	topic := ensureTopic(client, projectID, topicID)
	subscription := ensureSubscription(client, topic, subID)

	log.Printf("Have subscription: %v", subscription)
	//pullMsgs(projectID, topicID)
}
