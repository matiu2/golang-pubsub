/// This is the message that is sent across pubsub and eventually stored in the DB
package main

import (
	"math/rand"
	"time"
)

type Message struct {
	ServiceName string    `json:"service_name"`
	Payload     string    `json:"payload"`
	Severity    string    `json:"severity"`
	Timestamp   time.Time `json:"timestamp"`
}

/// Generates `n` randomish messages
func GenerateTestMessages(n int) []Message {
	severities := []string{"debug", "info", "warn", "error"}
	services := []string{"superman", "batman", "wonder woman", "captain kangaroo", "mister incredible"}
	payloads := []string{"flies", "jumps", "shoots", "smiles", "eats"}
	rand.Seed(time.Now().Unix())
	var out []Message
	for i := 0; i < n; i++ {
		msg := Message{
			ServiceName: services[rand.Intn(len(services))],
			Payload:     payloads[rand.Intn(len(payloads))],
			Severity:    severities[rand.Intn(len(severities))],
			Timestamp:   time.Now(),
		}
		out = append(out, msg)
	}
	return out
}
