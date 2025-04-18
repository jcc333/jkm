package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jcc333/jkm/internal/email"
	"github.com/jcc333/jkm/internal/log"
	"github.com/jcc333/jkm/internal/messages"
)

// Our custom commands for the application.
// These live in their own package versus by-model to avoid implementation-drift,
// commands being application-wide.

// ListView displays the list of messages.
func ListView() tea.Cmd {
	log.Info("list view command")

	return func() tea.Msg {
		log.Info("list message")
		return messages.ListMessages{}
	}
}

// ComposeView displays the compose message view.
func ComposeView() tea.Cmd {
	log.Info("compose view command")

	return func() tea.Msg {
		log.Info("compose message")
		return messages.ComposeMessage{}
	}
}

// ReadEmail displays the message details.
func ReadEmail(header email.MessageHeader) tea.Cmd {
	log.Info("read email command")

	return func() tea.Msg {
		log.Info("read email message")
		return messages.ReadEmailMessage{MessageHeader: &header}
	}
}

// SendEmail initiates the message sending process
func SendEmail(recipient, subject, body string) tea.Cmd {
	log.Info("send email command")

	return func() tea.Msg {
		log.Info("send email message")
		return messages.SendingMessage{
			Recipient: recipient,
			Subject:   subject,
			Body:      body,
		}
	}
}

// ShowError displays an error message
func ShowError(err error) tea.Cmd {
	log.Info("show error command")

	return func() tea.Msg {
		log.Info("show error message")
		return messages.Err{Error: err}
	}
}

// Fetch the message body for a message header
func FetchEmailBody(id int, receiver email.Receiver) tea.Cmd {
	log.Info("fetch email body command")

	return func() tea.Msg {
		body, err := receiver.Read(id)
		if err != nil {
			log.Errorf("error fetching body: %v", err)
			return messages.Err{Error: err}
		}
		log.Info("fetched body message")
		return messages.FetchedBody{ID: id, Body: body.Body}
	}
}
