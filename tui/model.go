package tui

import (
	"fmt"
	"strings"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/integrity"
	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/WarnetBes/cursor-tool/internal/storage"
	"github.com/WarnetBes/cursor-tool/internal/uuid"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateMenu state = iota
	stateConfirmReset
	stateResult
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).Padding(0, 1)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F59E0B"))
)

type Model struct {
	state    state
	cursor   int
	message  string
	isError  bool
	menuItems []string
}

func NewModel() Model {
	return Model{
		state: stateMenu,
		menuItems: []string{
			"Reset Machine ID",
			"Backup Storage",
			"Show Status",
			"Quit",
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.state == stateMenu && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == stateMenu && m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "enter":
			switch m.state {
			case stateMenu:
				switch m.cursor {
				case 0:
					m.state = stateConfirmReset
				case 1:
					m.message, m.isError = doBackup()
					m.state = stateResult
				case 2:
					m.message, m.isError = doStatus()
					m.state = stateResult
				case 3:
					return m, tea.Quit
				}
			case stateConfirmReset:
				m.message, m.isError = doReset()
				m.state = stateResult
			case stateResult:
				m.state = stateMenu
				m.message = ""
			}
		case "esc", "b":
			if m.state != stateMenu {
				m.state = stateMenu
				m.message = ""
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("cursor-tool") + "\n")
	sb.WriteString(mutedStyle.Render("Cursor IDE Machine ID Reset Utility") + "\n\n")

	switch m.state {
	case stateMenu:
		for i, item := range m.menuItems {
			if i == m.cursor {
				sb.WriteString(selectedStyle.Render("> "+item) + "\n")
			} else {
				sb.WriteString("  " + item + "\n")
			}
		}
		sb.WriteString("\n" + mutedStyle.Render("Use arrow keys to navigate, Enter to select, q to quit"))
	case stateConfirmReset:
		sb.WriteString("Are you sure you want to reset the Machine ID?\n\n")
		sb.WriteString("Press Enter to confirm, Esc to cancel.\n")
	case stateResult:
		if m.isError {
			sb.WriteString(errorStyle.Render("Error: "+m.message) + "\n")
		} else {
			sb.WriteString(successStyle.Render(m.message) + "\n")
		}
		sb.WriteString("\n" + mutedStyle.Render("Press Enter to go back"))
	}

	return sb.String()
}

func doReset() (string, bool) {
	paths, err := platform.GetPaths()
	if err != nil {
		return err.Error(), true
	}
	_ = backup.Create(paths)
	newIDs, err := uuid.GenerateIDs()
	if err != nil {
		return err.Error(), true
	}
	if err := storage.Write(paths.StorageJSON, newIDs); err != nil {
		return err.Error(), true
	}
	_ = integrity.WriteHMAC(paths.StorageJSON)
	_ = platform.WriteRegistry(newIDs)
	return fmt.Sprintf("Reset successful!\nmachineId: %s", newIDs.MachineID), false
}

func doBackup() (string, bool) {
	paths, err := platform.GetPaths()
	if err != nil {
		return err.Error(), true
	}
	path, err := backup.CreateWithPath(paths)
	if err != nil {
		return err.Error(), true
	}
	return "Backup created: " + path, false
}

func doStatus() (string, bool) {
	return "Status: use CLI 'status' command for full details", false
}
