package tui

import (
	"fmt"
	"strings"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/WarnetBes/cursor-tool/internal/storage"
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
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).Padding(0, 1)
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F59E0B"))
)

type Model struct {
	state     state
	cursor    int
	message   string
	isError   bool
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
		sb.WriteString("\n" + mutedStyle.Render("Arrow keys to navigate, Enter to select, q to quit"))
	case stateConfirmReset:
		sb.WriteString("Reset Machine ID? This cannot be undone.\n\n")
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
	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return err.Error(), true
	}
	mgr := backup.New(5)
	result, err := storage.ModifyStorageIDs(storagePath, mgr)
	if err != nil {
		return err.Error(), true
	}
	newID := ""
	for _, v := range result.After {
		newID = v
		break
	}
	return fmt.Sprintf("Reset successful!\nmachineId: %s", newID), false
}

func doBackup() (string, bool) {
	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return err.Error(), true
	}
	mgr := backup.New(5)
	path, err := mgr.Create(storagePath)
	if err != nil {
		return err.Error(), true
	}
	return "Backup created: " + path, false
}

func doStatus() (string, bool) {
	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return err.Error(), true
	}
	ids, err := storage.ReadCurrentIDs(storagePath)
	if err != nil {
		return err.Error(), true
	}
	var sb strings.Builder
	for _, k := range storage.TelemetryFields {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, ids[k]))
	}
	return sb.String(), false
}
