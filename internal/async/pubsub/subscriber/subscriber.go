package pubsub

import (
	"crypto/rand"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	async_models "github.com/danieeelfr/swd-challenge/internal/async/pubsub/models"
)

// Subscriber holds the subscriber implementation
type Subscriber struct {
	ID       string                     // id of subscriber
	Messages chan *async_models.Message // messages channel
	Topics   map[string]bool            // topics it is subscribed to.
	Active   bool                       // if given subscriber is active
	Mutex    sync.RWMutex               // lock
}

// CreateNewSubscriber returns a new subscriber implementation
func CreateNewSubscriber() (string, *Subscriber) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	id := fmt.Sprintf("%X-%X", b[0:4], b[4:8])

	return id, &Subscriber{
		ID:       id,
		Messages: make(chan *async_models.Message),
		Topics:   map[string]bool{},
		Active:   true,
	}
}

// AddTopic add a new topic
func (s *Subscriber) AddTopic(topic string) {
	// add topic to the subscriber
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	s.Topics[topic] = true
}

// Signal gets the message from the channel
func (s *Subscriber) Signal(msg *async_models.Message) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	if s.Active {
		s.Messages <- msg
	}
}

// Listen listens to the message channel, prints once received.
func (s *Subscriber) Listen() {
	for {
		if msg, ok := <-s.Messages; ok {
			log.Warn(msg.GetMessageBody())
		}
	}
}
