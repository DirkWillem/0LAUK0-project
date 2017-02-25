package dispatch

import "fmt"

type (
	// HTTPError contains error information for errors occurring in a dispatcher error
	DispatcherError struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}
)

func (err *DispatcherError) Error() string {
	return err.Message
}

// Creates a message based on the dispatcher error
func (err *DispatcherError) OutgoingMessage(requestID int) outgoingMessage {
	return outgoingMessage{
		SubscriptionID: -1,
		Action:         "error",
		RequestID:      requestID,
		Payload:        err,
	}
}

func ToDispatcherError(err error) *DispatcherError {
	if dpe, ok := err.(*DispatcherError); ok {
		return dpe
	}

	return &DispatcherError{
		Code:    "error",
		Message: err.Error(),
	}
}

// BadRequestErrorMessage creates a new bad request dispatcher error with the given message
func BadRequestErrorMessage(message string) *DispatcherError {
	return &DispatcherError{
		Code:    "bad_request",
		Message: message,
	}
}

// BadRequestError creates a new bad request dispatcher error with the given error
func BadRequestError(err error) *DispatcherError {
	return BadRequestErrorMessage(err.Error())
}

// UndefinedActionError creates a new undefined action dispatcher error with the given action
func UndefinedActionError(action string) *DispatcherError {
	return &DispatcherError{
		Code:    "undefined_action",
		Message: fmt.Sprintf("Undefined action %s", action),
	}
}

// UndefinedActionError creates a new undefined subject dispatcher error with the given subject
func UndefinedSubjectError(subject string) *DispatcherError {
	return &DispatcherError{
		Code:    "undefined_subject",
		Message: fmt.Sprintf("Undefined subject %s", subject),
	}
}

// UndefinedActionError creates a new already subscribed dispatcher error with the given subject
func AlreadySubscribedError(subject string) *DispatcherError {
	return &DispatcherError{
		Code:    "already_subscribed",
		Message: fmt.Sprintf("This client is already subscribed to subject %s with the same subscription parameters", subject),
	}
}