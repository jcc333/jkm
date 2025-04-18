package list

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jcc333/jkm/internal/email"
	"github.com/jcc333/jkm/internal/messages"
)

// Our Listing model.
//
// This component is responsible for browsing emails.
// That means, it is also responsible for receiving emails.
// It is also responsible for any sorting and searching functionality.

// The listing model owns
// - a reader
// - as composer
// - an error-handler

// Read/unread status of messages
type status int

const (
	unread status = iota
	read
)

type listingModel struct {
	messagesList  list.Model
	receiver      email.Receiver
	currentOffset int
	pageSize      int
	totalMessages int
	items         []*email.MessageHeader
}

// A list item for the `listingModel`.
// Contains an email header and implements `list.Item`.
type emailItem struct {
	header *email.MessageHeader
}

func (i emailItem) Title() string {
	return i.header.Subject
}

func (i emailItem) Description() string {
	return fmt.Sprintf("From: %s | %s", i.header.From, i.header.Date.Format(time.DateTime))
}

func (i emailItem) FilterValue() string {
	return i.header.Subject + " " + i.header.From
}

func New(receiver email.Receiver) *listingModel {
	delegate := list.NewDefaultDelegate()
	listModel := list.New([]list.Item{}, delegate, 0, 0)
	listModel.Title = "JKM Email Client"
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)

	pageSize := 50

	return &listingModel{
		messagesList:  listModel,
		receiver:      receiver,
		currentOffset: 0,
		pageSize:      pageSize,
		items:         make([]*email.MessageHeader, 0),
	}
}

func (m listingModel) Init() tea.Cmd {
	return tea.Batch(m.fetchMessages)
}

func (m listingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := 0, 0
		m.messagesList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case messages.FetchedList:
		// Replace items with the complete list
		m.items = msg.Items

		// Create list items for the UI
		listItems := make([]list.Item, len(m.items))
		for i, header := range m.items {
			listItems[i] = emailItem{header: header}
		}

		m.messagesList.SetItems(listItems)
		return m, tea.WindowSize()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.messagesList.SelectedItem().(emailItem)
			if ok {
				return m, tea.Batch(
					func() tea.Msg { return messages.ReadEmailMessage{MessageHeader: item.header} },
				)
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		case "c":
			return m, tea.Batch(
				func() tea.Msg { return messages.ComposeMessage{} },
			)

		case "r":
			m.items = make([]*email.MessageHeader, 0)
			return m, m.fetchMessages
		}
	}

	var cmd tea.Cmd
	m.messagesList, cmd = m.messagesList.Update(msg)
	return m, cmd
}

func (m listingModel) View() string {
	return m.messagesList.View()
}

// Fetch email overviews.
func (m listingModel) fetchMessages() tea.Msg {
	count, err := m.receiver.CountMessages()
	if err != nil {
		return messages.Err{Error: err}
	}
	m.totalMessages = count

	headers, err := m.receiver.List()
	if err != nil {
		return messages.Err{Error: err}
	}

	items := make([]*email.MessageHeader, len(headers))
	for i, header := range headers {
		headerCopy := header
		items[i] = &headerCopy
	}

	return messages.FetchedList{Items: items}
}
