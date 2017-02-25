package dispatch

import (
	"reflect"
)

type (
	// Subject contains the methods a subject struct must implement
	Subject interface {
		// GetTitle returns the title that can be used by clients to subscribe to the subject
		GetTitle() string

		// CreateSubscriptionParams creates a SubscriptionParams interface from a given set of subscription data
		CreateSubscriptionParams(subscriptionData map[string]interface{}) (SubscriptionParams, error)

		// ShouldSendMessageToSubscription returns whether a message should te sent to a subscription with certain subscription parameters
		MessageShouldBeSentToSubscription(message SubjectMessage, subscriptionParams SubscriptionParams) bool

		// GetMessageChan returns the outgoing message channel of the subject
		GetMessageChan() <-chan SubjectMessage
	}

	SubscriptionParams interface {
		// IsEqualTo returns whether the subscription params are equal to another SubscriptionParams
		IsEqualTo(subscriptionParams SubscriptionParams) bool
	}

	// subjectMessage contains the message data that is sent from a subject to the dispatcher
	SubjectMessage struct {
		Action  string
		Payload interface{}
	}

	// Dispatcher contains the information on a dispatcher
	Dispatcher struct {
		subjects []Subject
		clients  []*client
	}

	// Subscription contains information on a subscription to a subject of a client
	subscription struct {
		SubscriptionID     int
		SubjectTitle       string
		SubscriptionParams SubscriptionParams
	}

	// OutgoingMessage contains a message that is sent from a dispatcher to a client
	outgoingMessage struct {
		SubscriptionID int         `json:"subscriptionId"`
		Action         string      `json:"action"`
		RequestID      int         `json:"requestId"`
		Payload        interface{} `json:"payload"`
	}
)

// NewDispatcher creates a new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		subjects: make([]Subject, 0),
		clients:  make([]*client, 0),
	}
}

// RegisterSubject registers a new subject to a dispatcher
func (d *Dispatcher) RegisterSubject(subject Subject) {
	d.subjects = append(d.subjects, subject)
}

// CreateClient creates a new client in the dispatcher
func (d *Dispatcher) CreateClient() *client {
	client := newClient(d)

	d.clients = append(d.clients, client)

	return client
}

// Start starts the dispatcher process
func (d *Dispatcher) Start() {
	// Create a list of channels to which can be subscribed
	cases := make([]reflect.SelectCase, len(d.subjects))

	for i, subject := range d.subjects {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(subject.GetMessageChan()),
		}
	}

	// Keep on looping
	for true {
		// Read the next incoming message
		idx, value, _ := reflect.Select(cases)
		message := value.Interface().(SubjectMessage)
		subject := d.subjects[idx]

		// Iterate over all clients
		for _, c := range d.clients {

			// If the client is subscribed to the sending subject, send the message to the client
			if c.isSubscribedTo(subject.GetTitle()) {
				sub := c.getSubscription(subject.GetTitle())

				if subject.MessageShouldBeSentToSubscription(message, sub.SubscriptionParams) {
					c.OutgoingMessages <- outgoingMessage{
						SubscriptionID: sub.SubscriptionID,
						Action:         message.Action,
						Payload:        message.Payload,
						RequestID:      -1,
					}
				}
			}
		}
	}
}
