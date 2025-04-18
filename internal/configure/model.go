package configure

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/jcc333/jkm/internal/log"
	"github.com/jcc333/jkm/internal/messages"
)

// The configure model handles missing configuration details at the application's startup.
// That means, if it is initialized *with* a complete configuration, it will yield to the listing component.
// If it is not, it will prompt the user to complete the configuration and use that elsewhere.
type model struct {
	// The (global) configuration settings.
	cfg *Config

	// The form for setting our configuration.
	form *huh.Form

	// String carriers of the port ints.
	imapPort, smtpPort string
}

// Construct a new configuration model.
func New(cfg *Config) *model {
	var err error
	defer func() {
		if err != nil {
			log.Errorf("failed to build configure view '%v'", err)
		}
	}()

	log.Info("build configure view")
	m := &model{
		cfg:  cfg,
		form: nil,
	}
	fields := []huh.Field{
		huh.NewInput().Key("email").Title("Email Address").Value(&cfg.EmailAddress),
		huh.NewInput().Key("imapHost").Title("IMAP Server").Value(&cfg.IMAPServer),
		huh.NewInput().Key("imapPort").Title("IMAP Port").Value(&m.imapPort).Validate(func(s string) error {
			if _, err = strconv.Atoi(s); err != nil {
				return fmt.Errorf("invalid port number: %s", s)
			}
			return nil
		}),
		huh.NewInput().Key("imapPassword").Title("IMAP Password").EchoMode(huh.EchoModePassword).Value(&cfg.IMAPPassword),
		huh.NewInput().Key("smtpHost").Title("SMTP Server").Value(&cfg.SMTPServer),
		huh.NewInput().Key("smtpPort").Title("SMTP Port").Value(&m.smtpPort).Validate(func(s string) error {
			if _, err = strconv.Atoi(s); err != nil {
				return fmt.Errorf("invalid port number: %s", s)
			}
			return nil
		}),
		huh.NewInput().Key("imapPassword").Title("SMTP Password").EchoMode(huh.EchoModePassword).Value(&cfg.SMTPPassword),
	}

	form := huh.NewForm(
		huh.NewGroup(fields...).
			Title("Complete Configuration").
			Description("Missing configuration (add these to your .env to skip this). Press CTRL+C to quit."),
	)
	m.form = form
	return m
}

// The update for the configuring model mostly just handles the completed form.
// It is also responsible for forwarding the message to the form for handling.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debugf("update configure view: '%v'", msg)
	var err error
	defer func() {
		if err != nil {
			log.Errorf("failed to update configure view '%v'", err)
		}
	}()

	// Handle key interrupts (Ctrl+C)
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	// Pass the message to the form for handling
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	// If form has been completed, send a ConfigurationComplete message
	if m.form.State == huh.StateCompleted {
		m.cfg.IMAPPort, err = strconv.Atoi(m.imapPort)
		if err != nil {
			return m, cmd
		}
		m.cfg.SMTPPort, err = strconv.Atoi(m.smtpPort)
		if err != nil {
			return m, cmd
		}
		log.Info("configuration complete")
		return m, tea.Batch(cmd, func() tea.Msg { return messages.ListMessages{} })
	}

	return m, cmd
}

// Render the view.
func (m *model) View() string {
	return m.form.View()
}

// Init the underlying form.
func (m *model) Init() tea.Cmd {
	log.Info("init configure view")
	return m.form.Init()
}
