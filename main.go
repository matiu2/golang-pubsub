package main

import (
	"context"
	"database/sql"
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

	/// How many batches of messages to make/send total
	batches_to_send_orig := os.Getenv("BATCHES_TO_SEND")
	batches_to_send, err := strconv.Atoi(batches_to_send_orig)
	if err != nil {
		log.Fatalf("Unable to parse an int from env var: BATCHES_TO_SEND: %s, %v", batches_to_send_orig, err)
	}

	/// How many batches to store in memory
	buffer_count_orig := os.Getenv("BUFFER_COUNT")
	buffer_count, err := strconv.Atoi(buffer_count_orig)
	if err != nil {
		log.Fatalf("Unable to parse an int from env var: BUFFER_COUNT: %s, %v", buffer_count_orig, err)
	}

	/// How many msecs to delay between sending messages
	publish_delay_orig := os.Getenv("PUBLISH_DELAY_MSECS")
	publish_delay, err := strconv.Atoi(publish_delay_orig)
	if err != nil {
		log.Fatalf("Unable to parse an int from env var: PUBLISH_DELAY_MSECS: %s, %v", publish_delay_orig, err)
	}

	/// Make some batches of messages to send
	messages := GenerateTestMessages(batch_size * batches_to_send)

	// A channel for the pubsub puller to push to the database flusher
	// Allow up to 1000 batches before we stop accepting messages from pub sub
	batches := make(chan []Message, buffer_count)

	var wg sync.WaitGroup
	wg.Add(3)

	go PullMsgs(projectID, subID, len(messages), batch_size, batches, &wg)
	go PublishMsgs(projectID, topicID, messages, int64(publish_delay), &wg)
	go PeriodicFlush(batches, int64(flush_period), len(messages), &wg)

	log.Printf("Waiting for all the messages")
	wg.Wait()

	// Now check the database

	// Connect to the database
	log.Printf("Checking the database")
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/pubsub")
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	// See how many superman info messages there are
	rows, err := db.Query("Select count(1) from service_logs where service_name = ? and severity = ?", "superman", "debug")
	if err != nil {
		log.Fatalf("Unable to read supermans from DB: %v", err)
	}
	count := 0
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatalf("Unable to count supermans from DB: %v", err)
		}
	} else {
		log.Fatal("No supermans found in the DB")
	}

	// Make sure that matches up with the reports table
	rows, err = db.Query("select count from service_severity where service_name = ? and severity = ?", "superman", "debug")
	if err != nil {
		log.Fatalf("Unable to read superman summary from DB: %v", err)
	}
	new_count := 0
	if rows.Next() {
		err = rows.Scan(&new_count)
		if err != nil {
			log.Fatalf("Unable to check count supermans from DB: %v", err)
		}
	} else {
		log.Fatal("No supermans found in the report table in the DB")
	}

	if count != new_count {
		log.Fatalf("Wrong number supermans. Counted %d but status table had %d", count, new_count)
	} else {
		log.Printf("Everything worked as expected; We see %d superman debugs", count)
	}
}
