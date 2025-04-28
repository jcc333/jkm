package errorview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jcc333/jkm/internal/configure"
	"github.com/jcc333/jkm/internal/log"
	"github.com/jcc333/jkm/internal/messages"
)

// Model for handling and displaying errors
type model struct {
	cfg   *configure.Config
	error string
}

// New creates a new error model
func New(cfg *configure.Config, err error) *model {
	log.Info("new error model")
	return &model{
		cfg:   cfg,
		error: err.Error(),
	}
}

// Init initializes the error model
func (m *model) Init() tea.Cmd {
	log.Info("init error model")
	return nil
}

// Update handles updates to the error model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		log.Info("key interrupt error model")
		return m, tea.Batch(
			func() tea.Msg { return messages.ListMessages{} },
		)
	}
	return m, nil
}

// View renders the error model
func (m *model) View() string {
	log.Info("view error model")
	return fmt.Sprintf("Error: '%s'\n\nPress any key to return to the list view...", m.error)
}
