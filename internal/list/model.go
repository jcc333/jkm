package list

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jcc333/jkm/internal/email"
	"github.com/jcc333/jkm/internal/log"
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

// The model for a list of emails.
type listingModel struct {
	// The underlying list.
	list list.Model
	// If the list is currently being updated.
}

// A list item for the `listingModel`.
// Contains an email header and implements `list.Item`.
type emailItem struct {
	header *email.MessageHeader
}

// A list-item's title.
func (i emailItem) Title() string {
	return i.header.Subject
}

// A list-item's description.
func (i emailItem) Description() string {
	return fmt.Sprintf("From: %s | %s", i.header.From, i.header.Date.Format(time.DateTime))
}

// A list-item's search value.
func (i emailItem) FilterValue() string {
	return i.header.Subject + " " + i.header.From + " " + i.header.Date.Format(time.DateTime)
}

// Make a new mailer with the given index selected.
func New(items []*email.MessageHeader) *listingModel {
	delegate := list.NewDefaultDelegate()
	listModel := list.New([]list.Item{}, delegate, 0, 0)
	listModel.Title = "JKM Email Client"
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)

	listItems := make([]list.Item, len(items))
	for i, header := range items {
		listItems[i] = emailItem{header: header}
	}

	listModel.SetItems(listItems)

	return &listingModel{
		list: listModel,
	}
}

// Set up the list model.
func (m listingModel) Init() tea.Cmd {
	return nil
}

// List model update method.
func (m listingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := 0, 0
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case messages.RefreshedEmails:
		log.Warnf("Refreshing emails in list with %d messages", len(msg.Items))
		var selectedHeader *email.MessageHeader
		idx := m.list.Index()

		log.Warnf("idx = %v", idx)
		if idx >= 0 && idx < len(m.list.Items()) {
			if item, ok := m.list.SelectedItem().(emailItem); ok {
				selectedHeader = item.header
			}
		}
		log.Warnf("selected header = %v", selectedHeader)

		items := make([]list.Item, len(msg.Items))
		for i, header := range msg.Items {
			items[i] = emailItem{header: header}
		}
		log.Warnf("about to set items %d", len(items))

		m.list.SetItems(items)
		log.Warnf("set items %d", len(items))

		if selectedHeader != nil && len(items) > 0 {
			log.Warnf("have a selected header with Id %d and %d items", selectedHeader.ID, len(items))
			isFound := false

			for i, item := range items {
				if emailItem, ok := item.(emailItem); ok {
					if emailItem.header.ID == selectedHeader.ID {
						log.Warnf("selecting item %d with ID %d", i, emailItem.header.ID)
						m.list.Select(i)
						isFound = true
						break
					}
				}
			}

			if !isFound {
				if idx < len(items) {
					log.Warnf("falling back to item at same index")
					m.list.Select(idx)
				} else if len(items) > 0 {
					log.Warnf("falling back to last item")
					m.list.Select(len(items) - 1)
				}
			}
		} else if len(items) > 0 {
			// Ensure there's always a selection if there are items
			log.Warnf("falling back to first item")
			m.list.Select(0)
		}

		log.Warnf("returning model")
		return m, tea.WindowSize()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.list.SelectedItem().(emailItem)
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
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listingModel) View() string {
	return m.list.View()
}
