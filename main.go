package main

import (
	"context"
	"log"
	"os"
	"strconv"
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

	/// How many messages should we post to the database in one batch ?
	batch_size_orig := os.Getenv("BATCH_SIZE")
	batch_size, err := strconv.Atoi(batch_size_orig)
	if err != nil {
		log.Fatalf("Unable to parse an int from env var: BATCH_SIZE: %s, %v", batch_size_orig, err)
	}

	/// How often should we flush our messages (in seconds)
	flush_period_orig := os.Getenv("FLUSH_PERIOD")
	flush_period, err := strconv.Atoi(flush_period_orig)
	if err != nil {
		log.Fatalf("Unable to parse an int from env var: FLUSH_PERIOD: %s, %v", flush_period_orig, err)
	}

	/// How many batches to store in memory
	buffer_count_orig := os.Getenv("BUFFER_COUNT")
	buffer_count, err := strconv.Atoi(buffer_count_orig)
	if err != nil {
		log.Fatalf("Unable to parse an int from env var: BUFFER_COUNT: %s, %v", buffer_count_orig, err)
	}
	messages := GenerateTestMessages(batch_size * buffer_count)

	/// Make a 1000 batches of messages to send
	max := batch_size * 1000

	// A channel for the pubsub puller to push to the database flusher
	// Allow up to 1000 batches before we stop accepting messages from pub sub
	batches := make(chan []Message, 1000)
	var quit_flusher chan *sync.WaitGroup

	go PullMsgs(projectID, subID, max, batch_size, batches, &wg)
	go PeriodicPublish(projectID, topicID, messages, &wg)
	go PeriodicFlush(batches, int64(flush_period), quit_flusher)

	log.Printf("Waiting for all the messages")
	wg.Wait()

	// Send a quit to the Periodic flusher and wait for that too
	wg.Add(1)
	quit_flusher <- &wg
	wg.Wait()

}
