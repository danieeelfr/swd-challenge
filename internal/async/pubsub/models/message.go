package pubsub

// Message holds de message topic and body
type Message struct {
	Topic string
	Body  string
}

// NewMessage build and return a Message implementation
func NewMessage(msg string, topic string) *Message {
	return &Message{
		Topic: topic,
		Body:  msg,
	}
}

// GetMessageBody returns the message body
func (m *Message) GetMessageBody() string {
	return m.Body
}
