/// This is the message that is sent across pubsub and eventually stored in the DB
package main

import (
	"fmt"
	"time"
)

type Message struct {
	ServiceName string    `json:"service_name"`
	Payload     string    `json:"payload"`
	Severity    string    `json:"severity"`
	Timestamp   time.Time `json:"timestamp"`
}

/// Generates `n` random messages
func GenerateTestMessages(n int) []Message {
	var out []Message
	for i := 0; i < n; i++ {
		msg := Message{
			ServiceName: fmt.Sprintf("Service %d", i),
			Payload:     fmt.Sprintf("Payload %d", i),
			Severity:    fmt.Sprintf("Severity %d", i),
			Timestamp:   time.Now(),
		}
		out = append(out, msg)
	}
	return out
}
