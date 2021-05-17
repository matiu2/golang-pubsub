/// All database related work
package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/// Saves a batch of messages to the database
func SaveMessageBatch(messages []Message) {

	if len(messages) == 0 {
		return
	}

	log.Printf("Saving %d messages", len(messages))
	// Connect to the database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/pubsub")
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	// Start the database transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v", err)
	}

	// Counts service_name -> payload -> count added
	count_map := make(map[string]map[string]int)

	for _, msg := range messages {
		// Store the message stats for updating the reporting table later
		sub_map, ok := count_map[msg.ServiceName]
		if !ok {
			sub_map = map[string]int{}
			count_map[msg.ServiceName] = sub_map
		}
		_, ok = sub_map[msg.Severity]
		if !ok {
			sub_map[msg.Severity] = 1
		} else {
			sub_map[msg.Severity] += 1
		}

		// Insert each message into the database
		_, err := tx.Exec("INSERT INTO service_logs (service_name, payload, severity, timestamp) VALUES (?, ?, ?, ?)", msg.ServiceName, msg.Payload, msg.Severity, msg.Timestamp)
		if err != nil {
			log.Fatalf("Unable to insert message into database: %v", err)
		}
	}

	// Update the status table
	for service_name, severities := range count_map {
		for severity, count := range severities {
			_, err := tx.Exec("INSERT INTO service_severity (service_name, severity, count) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE count = count + ?", service_name, severity, count, count)
			if err != nil {
				log.Fatalf("Unable to update report: %v", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Unable to commit DB transaction: %v", err)
	}

	log.Printf("%d messages saved", len(messages))
}

/// Calls SaveMessages every minute
/// Every `period` seconds calls Save Messages - passing the batch channel, so it'll save all batches
/// Stops once it has saved 'max' messages
func PeriodicFlush(batches <-chan []Message, period int64, max int, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	count := 0
	for batch := range batches {
		// Wait for a bit
		<-ticker.C
		// Process the batch
		count += len(batch)
		SaveMessageBatch(batch)
		// Exit if we're done
		if count == max {
			break
		}
	}
}
