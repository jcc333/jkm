package commands

import (
	"time"

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
		return messages.SendingEmail{
			Recipient: recipient,
			Subject:   subject,
			Body:      body,
		}
	}
}

// Set the application to a sending-email spinner mode.
func SendingEmail(msg messages.SendEmail) tea.Cmd {
	log.Info("sending email command")

	return func() tea.Msg {
		log.Info("sending email")
		return messages.SendingEmail{
			Recipient: msg.Recipient,
			Subject:   msg.Subject,
			Body:      msg.Body,
		}
	}
}

// Handle a sent message.
func SentMessage() tea.Cmd {
	log.Info("sent email command")
	return func() tea.Msg {
		log.Info("sent email")
		return messages.SentEmail{}
	}
}

// A tick loop for updating our email listings.
func Tick() tea.Cmd {
	log.Info("tick command")
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		log.Info("tick")
		return messages.Tick(t)
	})
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

// RefreshEmails refreshes the list of messages.
func RefreshEmails(receiver email.Receiver, shouldBustCache bool) tea.Cmd {
	log.Info("refresh email command")

	return func() tea.Msg {
		log.Info("refresh emails")

		headers, err := receiver.List(shouldBustCache)
		if err != nil {
			return messages.Err{Error: err}
		}

		items := make([]*email.MessageHeader, len(headers))
		for i, header := range headers {
			headerCopy := header
			items[i] = &headerCopy
		}

		log.Info("refreshed emails")
		return messages.RefreshedEmails{Items: items}
	}
}
