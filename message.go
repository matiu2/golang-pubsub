/// This is the message that is sent across pubsub and eventually stored in the DB
package main

import "time"

type Message struct {
	ServiceName string    `json:"service_name"`
	Payload     string    `json:"payload"`
	Severity    string    `json:"severity"`
	Timestamp   time.Time `json:"timestamp"`
}
