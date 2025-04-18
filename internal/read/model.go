package read

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jcc333/jkm/internal/commands"
	"github.com/jcc333/jkm/internal/email"
	"github.com/jcc333/jkm/internal/messages"
)

// Our reading model.
// This component is responsible for displaying a single message.
type readingModel struct {
	// The underlying view for the pager.
	viewport viewport.Model

	// The keymap which governs the viewport controls.
	keyMap viewport.KeyMap

	// The message header being displayed
	header *email.MessageHeader

	// The full message including body (loaded after header)
	message *email.Message

	// Email client for fetching message details
	receiver email.Receiver
}

// Create a new reading model.
func New(receiver email.Receiver, header *email.MessageHeader) *readingModel {
	vp := viewport.New(0, 0)
	return &readingModel{
		viewport: vp,
		keyMap:   viewport.DefaultKeyMap(),
		header:   header,
		message:  nil,
		receiver: receiver,
	}
}

// Init the reading model as a tea view.
func (m readingModel) Init() tea.Cmd {
	if m.header != nil {
		return commands.FetchEmailBody(m.header.ID, m.receiver)
	}
	return nil
}

// Fetch the complete message
func (m readingModel) fetchFullMessage(id int) tea.Cmd {
	return func() tea.Msg {
		msg, err := m.receiver.Read(id)
		if err != nil {
			return messages.Err{Error: err}
		}
		return messages.FetchedOne{Message: msg}
	}
}

// Render the reading view.
func (m readingModel) View() string {
	headerStr := m.headerStatusView()

	var bodyStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		BorderTop(true).
		BorderBottom(true)

	return fmt.Sprintf("%s\n%s",
		headerStr,
		bodyStyle.Render(m.viewport.View()))
}

func (m readingModel) headerView() string {
	var headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	if m.header == nil {
		return headerStyle.Render("No message selected")
	}

	from := fmt.Sprintf("From: %s", m.header.From)
	to := fmt.Sprintf("To: %s", strings.Join(m.header.To, ", "))
	subject := fmt.Sprintf("Subject: %s", m.header.Subject)
	date := fmt.Sprintf("Date: %s", m.header.Date.Format(time.DateTime))

	return headerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			subject,
			from,
			to,
			date,
		),
	)
}

// New compact header view as a status bar
func (m readingModel) headerStatusView() string {
	if m.header == nil {
		return "No message selected"
	}

	var headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		BorderBottom(true).
		Width(100).
		Bold(true)

	// Show a compact header with just subject and from
	return headerStyle.Render(fmt.Sprintf("From: %s | Subject: %s",
		m.header.From, m.header.Subject))
}

func (m readingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, commands.ListView()
		case "j", "down":
			m.viewport.ScrollDown(1)
		case "k", "up":
			m.viewport.ScrollUp(1)
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerStatusView())
		contentHeight := msg.Height - headerHeight - 2
		m.viewport.Width = msg.Width
		m.viewport.Height = contentHeight
		if m.message != nil {
			m.viewport.SetContent(m.message.Body)
		} else if m.header != nil {
			m.viewport.SetContent("Loading message body...")
		}

	case messages.FetchedBody:
		// We got just the body content - create a full message from our header and this body
		if m.header != nil && m.header.ID == msg.ID {
			m.message = &email.Message{
				MessageHeader: *m.header,
				Body:          msg.Body,
			}
			if msg.Body == "" {
				m.viewport.SetContent("[No message body available]")
			} else {
				m.viewport.SetContent(msg.Body)
			}
			cmds = append(cmds, tea.WindowSize())
		}

	case messages.FetchedOne:
		m.message = msg.Message
		if m.message != nil {
			m.header = &m.message.MessageHeader
			if m.message.Body == "" {
				m.viewport.SetContent("[No message body available]")
			} else {
				m.viewport.SetContent(m.message.Body)
			}
		} else {
			m.viewport.SetContent("Error: Message could not be fetched")
		}
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
