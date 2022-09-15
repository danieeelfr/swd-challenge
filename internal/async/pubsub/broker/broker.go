package pubsub

import (
	"fmt"
	"sync"
	"time"

	async_models "github.com/danieeelfr/swd-challenge/internal/async/pubsub/models"
	subscriber "github.com/danieeelfr/swd-challenge/internal/async/pubsub/subscriber"
)

// Subscribers holds the subscribers
type Subscribers map[string]*subscriber.Subscriber

// Broker holds the Broker implementation
type Broker struct {
	subscribers Subscribers            // map of subscribers id:Subscriber
	topics      map[string]Subscribers // map of topic to subscribers
	mut         sync.RWMutex           // mutex lock
}

// NewBroker returns Broker object
func NewBroker() *Broker {
	return &Broker{
		subscribers: Subscribers{},
		topics:      map[string]Subscribers{},
	}
}

// AddSubscriber Add subscriber to the broker.
func (b *Broker) AddSubscriber() *subscriber.Subscriber {
	b.mut.Lock()
	defer b.mut.Unlock()
	id, s := subscriber.CreateNewSubscriber()
	b.subscribers[id] = s
	return s
}

// Subscribe subscribe to given topic
func (b *Broker) Subscribe(s *subscriber.Subscriber, topic string) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}
	s.AddTopic(topic)
	b.topics[topic][s.ID] = s
	fmt.Printf("%s Subscribed for topic: %s\n", s.ID, topic)
}

// Publish the given message to the also given topic.
func (b *Broker) Publish(topic string, msg string) {
	b.mut.RLock()
	bTopics := b.topics[topic]

	// TODO: it's here only to show that it is not blocking the main thread
	time.Sleep(time.Second * 3)

	b.mut.RUnlock()
	for _, s := range bTopics {
		m := async_models.NewMessage(msg, topic)
		if !s.Active {
			return
		}
		go (func(s *subscriber.Subscriber) {
			s.Signal(m)
		})(s)
	}
}
