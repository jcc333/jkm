package messages

import (
	"github.com/jcc333/jkm/internal/email"
)

// Our application's custom message types.

// The result of asynchronously fetching a page of email headers.
type FetchedList struct {
	Items []*email.MessageHeader
}

// The result of asynchronously fetching a single complete email.
type FetchedOne struct {
	Message *email.Message
}

// FetchedBody represents the result of fetching just the body content of a message
type FetchedBody struct {
	ID   int
	Body string
}

// An envelope for an error
type Err struct {
	Error error
}

// ComposeMessage represents the user composing a message.
type ComposeMessage struct{}

// ListMessages is sent when it's time to show the list view.
type ListMessages struct{}

// ReadEmailMessage is sent when a message is selected to be read.
type ReadEmailMessage struct {
	MessageHeader *email.MessageHeader
}

// Sent when we send a message.
type SendMessage struct {
	Recipient, Subject, Body string
}

// Sent when we *have sent* a message.
type SentMessage struct{}

// LoadMore is sent when the user has scrolled to the bottom of the list
// and we need to load more messages
type LoadMore struct {
	Offset int
}

// SendingMessage is sent when we are in the process of sending a message
// This triggers showing a spinner overlay
type SendingMessage struct {
	Recipient, Subject, Body string
}

// SendingFailure is sent when sending a message failed
type SendingFailure struct {
	Error error
}
