package publisher

import (
	"math/rand"

	broker "github.com/danieeelfr/swd-challenge/internal/async/pubsub/broker"
)

// available topics
var availableTopics = map[string]string{
	"NOTIFY": "MANAGER",
}

// Notify notify the subscribers
func Notify(broker *broker.Broker, message string) {
	topicValues := make([]string, 0, len(availableTopics))
	for _, v := range availableTopics {
		topicValues = append(topicValues, v)
	}

	v := topicValues[rand.Intn(len(topicValues))] // all topic values.
	go broker.Publish(v, message)
}
