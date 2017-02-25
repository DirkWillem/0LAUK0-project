package dispatch

import (
	"fmt"
	"reflect"
)

type (
	// Client contains information on a client of the dispatcher
	client struct {
		Dispatcher       *Dispatcher
		Subscriptions    []subscription
		OutgoingMessages chan outgoingMessage
	}

	// Contains information on an incoming message
	incomingMessage struct {
		Action    string                 `json:"action"`
		RequestID int                    `json:"requestId"`
		Payload   map[string]interface{} `json:"payload"`
	}

	// Payload for a subscription response
	subscriptionPayload struct {
		SubscriptionID int `json:"subscriptionId"`
	}
)

// newClient creates a new Client
func newClient(dispatcher *Dispatcher) *client {
	return &client{
		Dispatcher:       dispatcher,
		Subscriptions:    make([]subscription, 0),
		OutgoingMessages: make(chan outgoingMessage, 10),
	}
}

// isSubscribedTo returns whether a client is subscribed to a subject with the given title
func (c *client) isSubscribedTo(subjectTitle string) bool {
	for _, subscription := range c.Subscriptions {
		if subscription.SubjectTitle == subjectTitle {
			return true
		}
	}

	return false
}

// getSubscription returns a subscription for a given subject title
func (c *client) getSubscription(subjectTitle string) subscription {
	for _, subscription := range c.Subscriptions {
		if subscription.SubjectTitle == subjectTitle {
			return subscription
		}
	}

	return subscription{}
}

// Subscribe subscribes a client to a subscription
func (c *client) Subscribe(subjectTitle string, subscriptionParams map[string]interface{}) (subscription, error) {
	for _, subject := range c.Dispatcher.subjects {
		params := subject.CreateSubscriptionParams(subscriptionParams)

		// Check whether the user is already subscribed
		for _, sub := range c.Subscriptions {
			if sub.SubjectTitle == subjectTitle {
				if sub.SubscriptionParams.IsEqualTo(params) {
					return subscription{}, AlreadySubscribedError(subjectTitle)
				}
			}
		}

		// Subscribe
		if title := subject.GetTitle(); title == subjectTitle {
			c.Subscriptions = append(c.Subscriptions, subscription{
				SubscriptionID:     len(c.Subscriptions),
				SubjectTitle:       title,
				SubscriptionParams: params,
			})

			return c.Subscriptions[len(c.Subscriptions)-1], nil
		}
	}

	return subscription{}, UndefinedSubjectError(subjectTitle)
}

// Unsubscribe unsubscribes a client from a subscription with the given subscription ID
func (c *client) Unsubscribe(subscriptionID int) {
	c.Subscriptions[subscriptionID] = c.Subscriptions[len(c.Subscriptions)-1]
	c.Subscriptions = c.Subscriptions[:len(c.Subscriptions)-1]
}

// handleIncomingMessage handles an incoming message from the client
func (c *client) handleIncomingMessage(msg incomingMessage) {
	switch msg.Action {
	case "subscribe":
		// Type checks
		s, ok := msg.Payload["subject"]
		if !ok {
			c.OutgoingMessages <- BadRequestErrorMessage("Missing field 'subject' in payload of subscription action").OutgoingMessage(msg.RequestID)
			return
		}

		subjectTitle, ok := s.(string)
		if !ok {
			c.OutgoingMessages <- BadRequestErrorMessage(fmt.Sprintf("Invalid type for field 'subject' in payload of subscription action: expected string, got %s", reflect.TypeOf(s).Name())).OutgoingMessage(msg.RequestID)
			return
		}

		sp, ok := msg.Payload["subscriptionParams"]
		if !ok {
			c.OutgoingMessages <- BadRequestErrorMessage("Missing field 'subject' in payload of subscription action").OutgoingMessage(msg.RequestID)
			return
		}

		subscriptionParams, ok := sp.(map[string]interface{})
		if !ok {
			c.OutgoingMessages <- BadRequestErrorMessage(fmt.Sprintf("Invalid type for field 'subject' in payload of subscription action: expected object, got %s", reflect.TypeOf(sp).Name())).OutgoingMessage(msg.RequestID)
			return
		}

		// Subscribe
		subscription, err := c.Subscribe(subjectTitle, subscriptionParams)

		if err != nil {
			c.OutgoingMessages <- ToDispatcherError(err).OutgoingMessage(msg.RequestID)
			return
		}

		c.OutgoingMessages <- outgoingMessage{
			SubscriptionID: -1,
			RequestID: msg.RequestID,
			Action: "subscribe",
			Payload: subscriptionPayload{SubscriptionID: subscription.SubscriptionID},
		}

	default:
		c.OutgoingMessages <- UndefinedActionError(msg.Action).OutgoingMessage(msg.RequestID)
	}

}
