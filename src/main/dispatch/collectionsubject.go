package dispatch

type (
	// CollectionSubject represents a subscribable subject pertaining to a collection of entities
	CollectionSubject struct {
		Title    string
		messages chan SubjectMessage
	}

	// collectionSubjectSubscriptionParams contains the subscription parameters to a CollectionSubject
	collectionSubjectSubscriptionParams struct {
	}

	// CollectionEntityAddedPayload contains the payload for an "added" message
	CollectionEntityAddedPayload struct {
		ID          int         `json:"id"`
		AddedEntity interface{} `json:"addedEntity"`
	}

	// CollectionEntityUpdatedPayload contains the payload for an "updated" message
	CollectionEntityUpdatedPayload struct {
		ID            int         `json:"id"`
		UpdatedEntity interface{} `json:"updatedEntity"`
	}

	// CollectionEntityDeletedPayload contains the payload for a "deleted" message
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
		messages: make(chan SubjectMessage, 10),
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

func (cs *CollectionSubject) CreateSubscriptionParams(params map[string]interface{}) (SubscriptionParams, error) {
	return &collectionSubjectSubscriptionParams{}, nil
}

func (cs *CollectionSubject) MessageShouldBeSentToSubscription(message SubjectMessage, subscriptionParams SubscriptionParams) bool {
	return true
}

func (cs *CollectionSubject) GetMessageChan() <-chan SubjectMessage {
	return cs.messages
}

// EntityAdded notifies subscribers of the subject that a new entity has been added
func (cs *CollectionSubject) EntityAdded(entityID int, addedEntity interface{}) {
	cs.messages <- SubjectMessage{
		Action: CollectionEntityAddedAction,
		Payload: CollectionEntityAddedPayload{
			ID:          entityID,
			AddedEntity: addedEntity,
		},
	}
}

// EntityUpdated notifies subscribers of the subject that an entity has been updated
func (cs *CollectionSubject) EntityUpdated(entityID int, updatedEntity interface{}) {
	cs.messages <- SubjectMessage{
		Action: CollectionEntityUpdatedAction,
		Payload: CollectionEntityUpdatedPayload{
			ID:            entityID,
			UpdatedEntity: updatedEntity,
		},
	}
}

// EntityDeleted notifies subscribers of the subject that an entity has been deleted
func (cs *CollectionSubject) EntityDeleted(entityID int) {
	cs.messages <- SubjectMessage{
		Action: CollectionEntityDeletedAction,
		Payload: CollectionEntityUpdatedPayload{
			ID: entityID,
		},
	}
}
