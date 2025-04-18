package compose

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/jcc333/jkm/internal/commands"
	"github.com/jcc333/jkm/internal/configure"
	"github.com/jcc333/jkm/internal/log"
	"github.com/jcc333/jkm/internal/messages"
)

// Our model for composing emails.
// This is a pretty trivial huh form view.

// Our composing model has 3 fields (recipient, subject, body)
// This type enumerates them to make focus easier.
type field int

const (
	recipient field = iota
	subject
	body
)

// Our composing model.
// Has text inputs for the subject, recipient, and body.
// Has a pointer to the global config.
// Has a pointer to the field being edited.
type model struct {
	// The global configuration.
	cfg *configure.Config

	// The form for composing our email.
	form *huh.Form

	// The recipient(s), subject, and body of the email
	recipient, subject, body string

	// Whether the email has been confirmed for sending.
	isConfirmed bool
}

// Construct a new composing model.
func New(cfg *configure.Config) *model {
	log.Info("compose: initializing compose model")

	var m model
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("recipient").Title("Recipient").Value(&m.recipient),
			huh.NewInput().Key("subject").Title("Subject").Value(&m.subject),
			huh.NewText().Key("body").Title("Body").Value(&m.body),
			huh.NewConfirm().
				Title("Send an Email?").
				Affirmative("Send").
				Negative("Cancel").
				Value(&m.isConfirmed),
		),
	)
	m.form = form
	return &m
}

// The update for the composing model handles the flows for discarding a draft or else sending it.
// It is also responsible for forwarding the message to the form for handling.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type.String() {
		case "ctrl+c", "q", "esc":
			log.Debug("compose: canceling compose and returning to list view")
			return m, commands.ListView()
		}
	}

	switch msg.(type) {
	case messages.SentMessage:
		log.Info("compose: email sent successfully, returning to list view")
		return m, commands.ListView()
	}

	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		if m.isConfirmed {
			log.Info(fmt.Sprintf("compose: sending to %s, subject: %s", m.recipient, m.subject))
			return m, commands.SendEmail(m.recipient, m.subject, m.body)
		}
		if !m.isConfirmed {
			log.Debug("compose: user canceled sending, returning to list view")
			return m, commands.ListView()
		}
	}

	return m, cmd
}

// Render the view.
func (m *model) View() string {
	return m.form.View()
}

// Init the underlying form.
func (m *model) Init() tea.Cmd {
	log.Info("compose: initializing compose form")
	m.form.Init()
	return textinput.Blink
}
