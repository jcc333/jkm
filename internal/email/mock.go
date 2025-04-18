package email

import (
	"time"
)

// Mock values for the email package for testing.

// Mock is a SenderReceiver.
type Mock struct {
	// Inbound emails.
	inbox []MessageHeader

	// Outbound emails.
	outbox []Message
}

func NewMock() *Mock {
	return &Mock{
		inbox: []MessageHeader{
			{
				ID:      1,
				From:    "alice@example.com",
				To:      []string{"bob@example.com"},
				Subject: "Hello World",
				Date:    time.Now().Add(-time.Hour * 24),
			},
			{
				ID:      2,
				From:    "carol@example.com",
				To:      []string{"alice@example.com", "bob@example.com"},
				Subject: "Meeting Tomorrow",
				Date:    time.Now().Add(-time.Hour * 12),
			},
			{
				ID:      3,
				From:    "bob@example.com",
				To:      []string{"alice@example.com"},
				Subject: "Re: Hello World",
				Date:    time.Now().Add(-time.Hour * 6),
			},
		},
		outbox: []Message{},
	}
}

// Disconnect from the server: NOOP
func (m *Mock) Disconnect() error {
	return nil
}

// Send an email: NOOP
func (m *Mock) Send(msg Message) error {
	m.outbox = append(m.outbox, msg)
	return nil
}

// List some fake emails.
func (m *Mock) List(limit, offset int) ([]MessageHeader, error) {
	return m.inbox, nil
}

// Read a fake email.
func (m *Mock) Read(id int) (*Message, error) {
	for _, header := range m.inbox {
		if header.ID == id {
			return &Message{
				MessageHeader: header,
				Body:          "This is a very real email.\nIt is also a very good email.\n\nAll the best,\n\t-Sender",
			}, nil
		}
	}
	return nil, nil
}

// The number of messages in the inbox.
func (m *Mock) CountMessages() (int, error) {
	return len(m.inbox), nil
}
