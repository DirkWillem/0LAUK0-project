package dispatch

type (
	// CollectionSubject represents a subscribable subject pertaining to a single entity
	EntitySubject struct {
		Title    string
		messages chan subjectMessage
	}

	// entitySubjectSubscriptionParams contains the subscription parameters to a EntitySubject
	entitySubjectSubscriptionParams struct {
		EntityID int
	}

	entityUpdatedPayload struct {

	}
)
