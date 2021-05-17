/// pubsub message puller

package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

/// Pull `max` messages, push them into the DB via the channel, then return
func PullMsgs(projectID, subID string, max int, batch_size int, out chan<- []Message, done *sync.WaitGroup) {
	defer done.Done()
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	log.Printf("Listening for messages: %s - %s", projectID, subID)

	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	var batch []Message
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		log.Printf("Got message: %q\n", string(msg.Data))
		// Turn the msg data into a struct again
		var incoming Message
		err = json.Unmarshal(msg.Data, &incoming)
		if err != nil {
			log.Fatalf("Unable to demarshal a received message: %v", err)
		}
		msg.Ack()
		received++
		// Store the message in the batch
		batch = append(batch, incoming)

		// Once we've reached the batch size, push them out the channel
		if len(batch) == batch_size {
			log.Printf("Sending a batch of messages to the database")
			out <- batch
			batch = batch[:0]
		}

		// Once we've received the max messages, fire off one more batch then exit
		if received == max {
			log.Printf("Maximum messages (%d) have bee pulled", max)
			out <- batch
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Receive: %v", err)
	}

	log.Print("PULLER EXITING")
}
