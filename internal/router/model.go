package router

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jcc333/jkm/internal/commands"
	"github.com/jcc333/jkm/internal/compose"
	"github.com/jcc333/jkm/internal/configure"
	"github.com/jcc333/jkm/internal/email"
	"github.com/jcc333/jkm/internal/errorview"
	"github.com/jcc333/jkm/internal/io"
	"github.com/jcc333/jkm/internal/list"
	"github.com/jcc333/jkm/internal/log"
	"github.com/jcc333/jkm/internal/messages"
	"github.com/jcc333/jkm/internal/read"
	"github.com/jcc333/jkm/internal/sending"
)

type mode int

const (
	// Configuring the application in absence of a full .env file or environment variables.
	configureMode mode = iota

	// Listing the emails of the user.
	listMode

	// Reading an email
	readMode

	// Writing an email
	composeMode

	// Sending an email
	sendingMode

	// Recovering from an error
	errorMode
)

// The router model handles top-level events, and determines the member model which will View and Update.
type model struct {
	// The current mode of the router.
	mode mode

	// The previous mode of the router.
	previous *mode

	// The current model for the router's mode.
	model tea.Model

	// The (global) configuration settings.
	cfg *configure.Config

	// The email layer underpinning the application.
	mailer email.Client

	// Track if we're currently sending an email to prevent duplicates, continue the spinner.
	isSending bool
}

// Build a new router model.
func New(cfg *configure.Config) (*model, error) {
	log.Info("build router")

	// Determine initial mode based on configuration completeness
	initialMode := configureMode
	if cfg.IMAPServer != "" && cfg.EmailAddress != "" && cfg.IMAPPassword != "" {
		initialMode = listMode
	}
	var mailer email.Client

	m := &model{
		mode:      initialMode,
		cfg:       cfg,
		mailer:    mailer,
		isSending: false,
	}

	if initialMode == configureMode {
		log.Info("starting in configure mode")
		m.model = configure.New(cfg)
	} else {
		log.Info("starting in list mode")
		err := m.buildMailer()
		if err != nil {
			return nil, err
		}
		m.model = list.New([]*email.MessageHeader{})
	}

	return m, nil
}

// Build the mailer for the application.
func (m *model) buildMailer() error {
	log.Info("build mailer")
	if m.mailer != nil {
		return nil
	}
	mailer, err := io.New(m.cfg)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	m.mailer = mailer
	return nil
}

// Close the model's mailer
func (m *model) Disconnect() error {
	log.Info("disconnecting mailer")
	return m.mailer.Disconnect()
}

// Initialize the router model.
func (m *model) Init() tea.Cmd {
	log.Info("init router")
	return tea.Sequence(commands.RefreshEmails(m.mailer, true), m.model.Init(), commands.Tick())
}

// Handle updates to the model, generally by routing them to the routed model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.mode == listMode {
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
		}
	}

	// Check for mode switch messages first
	switch msg := msg.(type) {
	case messages.ListMessages:
		return m, tea.Sequence(m.list(), commands.RefreshEmails(m.mailer, false))

	case messages.ReadEmailMessage:
		return m, tea.Sequence(m.read(msg.MessageHeader), commands.FetchEmailBody(msg.MessageHeader.ID, m.mailer))

	case messages.Err:
		return m, m.recover(msg.Error)

	case messages.SendingFailure:
		m.isSending = false
		return m, m.recover(msg.Error)

	case messages.ComposeMessage:
		return m, m.compose()

	case messages.SendEmail:
		return m, commands.SendingEmail(msg)

	case messages.SendingEmail:
		if m.isSending {
			return m, nil
		}

		m.isSending = true

		sendingCmd := m.sending(msg.Recipient, msg.Subject, msg.Body)
		sendCmd := m.sendMessage(msg.Recipient, msg.Subject, msg.Body)
		return m, tea.Batch(sendingCmd, sendCmd)

	case messages.SentEmail:
		m.isSending = false
		if m.mode == composeMode {
			return m, commands.SentMessage()
		} else {
			return m, tea.Sequence(m.list(), commands.RefreshEmails(m.mailer, true))
		}

	case messages.Tick:
		return m, tea.Batch(commands.RefreshEmails(m.mailer, false), commands.Tick())
	}

	model, cmd := m.model.Update(msg)
	m.model = model

	return m, cmd
}

func (m *model) sendMessage(recipient string, subject string, body string) tea.Cmd {
	return func() tea.Msg {
		msg := email.Message{
			MessageHeader: email.MessageHeader{
				From:    m.cfg.EmailAddress,
				To:      []string{recipient},
				Subject: subject,
			},
			Body: body,
		}

		err := m.mailer.Send(msg)
		if err != nil {
			return messages.SendingFailure{Error: err}
		}

		return messages.SentEmail{}
	}
}

// Render the view - delegates to the current routed model's View method.
func (m *model) View() string {
	return m.model.View()
}

// List the emails in the inbox.
func (m *model) list() tea.Cmd {
	log.Info("listing emails")
	if m.mailer == nil {
		err := m.buildMailer()
		if err != nil {
			fmt.Printf("Error building mailer: %v\n", err)
			os.Exit(1)
		}
	}
	m.mode = listMode
	m.model = list.New([]*email.MessageHeader{})
	return m.model.Init()
}

// Read an email.
func (m *model) read(header *email.MessageHeader) tea.Cmd {
	m.mode = readMode
	m.model = read.New(m.mailer, header)
	return m.model.Init()
}

// Recover from an error.
func (m *model) recover(err error) tea.Cmd {
	m.mode = errorMode
	m.model = errorview.New(m.cfg, err)
	return m.model.Init()
}

// Compose an email.
func (m *model) compose() tea.Cmd {
	m.mode = composeMode
	m.model = compose.New(m.cfg)
	return m.model.Init()
}

// Wait while the email sends.
func (m *model) sending(recipient, subject, body string) tea.Cmd {
	m.mode = sendingMode
	m.model = sending.New(recipient, subject, body)
	return m.model.Init()
}
