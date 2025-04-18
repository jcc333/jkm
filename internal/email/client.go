package email

import (
	"time"
)

// MessageHeader represents the header information of an email message
// Used for efficient listing without loading full message content
type MessageHeader struct {
	ID      int
	From    string
	To      []string
	Subject string
	Date    time.Time
}

// Message represents a complete email message including body content
type Message struct {
	MessageHeader
	Body string
}

// A type for sending emails.
type Sender interface {
	Send(msg Message) error
}

// A type for receiving emails.
type Receiver interface {
	List() ([]MessageHeader, error)
	Read(id int) (*Message, error)
	CountMessages() (int, error)
}

// A type for sending or receiving emails.
type Client interface {
	Sender
	Receiver
	Disconnect() error
}
