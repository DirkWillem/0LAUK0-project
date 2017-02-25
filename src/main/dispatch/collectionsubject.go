package dispatch

type (
	// CollectionSubject represents a subscribable subject pertaining to a collection of entities
	CollectionSubject struct {
		Title    string
		messages chan subjectMessage
	}

	// collectionSubjectSubscriptionParams contains the subscription parameters to a CollectionSubject
	collectionSubjectSubscriptionParams struct {
	}

	CollectionEntityAddedPayload struct {
		ID          int         `json:"id"`
		AddedEntity interface{} `json:"addedEntity"`
	}

	CollectionEntityUpdatedPayload struct {
		ID            int         `json:"id"`
		UpdatedEntity interface{} `json:"updatedEntity"`
	}

	CollectionEntityDeletedPayload struct {
		ID int `json:"id"`
	}
)

const (
	CollectionEntityAddedAction   = "added"
	CollectionEntityUpdatedAction = "updated"
	CollectionEntityDeletedAction = "deleted"
)

// NewCollectionSubject creates a new CollectionSubject
func NewCollectionSubject(title string, dispatcher *Dispatcher) *CollectionSubject {
	subject := &CollectionSubject{
		Title:    title,
		messages: make(chan subjectMessage, 10),
	}

	dispatcher.RegisterSubject(subject)

	return subject
}

func (cssp *collectionSubjectSubscriptionParams) IsEqualTo(params SubscriptionParams) bool {
	return true
}

func (cs *CollectionSubject) GetTitle() string {
	return cs.Title
}

func (cs *CollectionSubject) CreateSubscriptionParams(params map[string]interface{}) SubscriptionParams {
	return &collectionSubjectSubscriptionParams{}
}

func (cs *CollectionSubject) MessageShouldBeSentToSubscription(message subjectMessage, subscriptionParams SubscriptionParams) bool {
	return true
}

func (cs *CollectionSubject) GetMessageChan() <-chan subjectMessage {
	return cs.messages
}

// Notifies subscribers of the subject that a new entity has been added
func (cs *CollectionSubject) EntityAdded(entityID int, addedEntity interface{}) {
	cs.messages <- subjectMessage{
		Action: CollectionEntityAddedAction,
		Payload: CollectionEntityAddedPayload{
			ID:          entityID,
			AddedEntity: addedEntity,
		},
	}
}

// Notifies subscribers of the subject that an entity has been updated
func (cs *CollectionSubject) EntityUpdated(entityID int, updatedEntity interface{}) {
	cs.messages <- subjectMessage{
		Action: CollectionEntityUpdatedAction,
		Payload: CollectionEntityUpdatedPayload{
			ID:            entityID,
			UpdatedEntity: updatedEntity,
		},
	}
}

// Notifies subscribers of the subject that an entity has been deleted
func (cs *CollectionSubject) EntityDeleted(entityID int) {
	cs.messages <- subjectMessage{
		Action: CollectionEntityDeletedAction,
		Payload: CollectionEntityUpdatedPayload{
			ID: entityID,
		},
	}
}
