package sending

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jcc333/jkm/internal/commands"
	"github.com/jcc333/jkm/internal/messages"
)

// model for the sending state view.
type model struct {
	// Sending state management
	spinner spinner.Model

	// Email details
	recipient, subject, body string
}

// New creates a new sending model.
func New(recipient, subject, body string) *model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &model{
		spinner:   s,
		recipient: recipient,
		subject:   subject,
		body:      body,
	}
}

// Update handles messages for the sending model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.SentEmail:
		// Email sent successfully
		return m, commands.ListView()
	case messages.SendingFailure:
		// Email sending failed
		return m, commands.ShowError(msg.Error)
	}

	// Update the spinner
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the sending overlay.
func (m *model) View() string {
	// Create sending overlay
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 3).
		Align(lipgloss.Center).
		Bold(true)

	return style.Render(m.spinner.View() + " Sending email...")
}

// Init initializes the sending model.
func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}
