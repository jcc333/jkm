package messages

import (
	"time"

	"github.com/jcc333/jkm/internal/email"
)

// Our application's custom message types.

// The result of asynchronously fetching a page of email headers.
type RefreshedEmails struct {
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
type SendEmail struct {
	Recipient, Subject, Body string
}

// Sent when we *have sent* a message.
type SentEmail struct{}

// SendingEmail is sent when we are in the process of sending a message
// This triggers showing a spinner overlay
type SendingEmail struct {
	Recipient, Subject, Body string
}

// SendingFailure is sent when sending a message failed
type SendingFailure struct {
	Error error
}

// A tick event. Used in our case to refresh the email list.
type Tick time.Time
