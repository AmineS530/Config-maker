package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ClearTerminal clears the terminal screen cleanly using ANSI escape sequences.
// This is used to ensure a fresh, clutter-free terminal viewport before and after menu displays.
func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

// menuModel represents the state of the interactive Bubble Tea selection menu.
type menuModel struct {
	cursor   int      // The index of the item currently pointed to by the cursor (0-indexed).
	choices  []string // The list of text choices/options presented in the menu.
	selected int      // The index of the selected item, or -1 if the user exited without choosing.
}

// Init initializes the Bubble Tea program. Since this model has no background
// tasks or active setups at start, it returns nil.
func (m menuModel) Init() tea.Cmd {
	return nil
}

// Update handles state changes inside the main menu loop based on incoming key events.
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Exit keys: quit the menu cleanly and mark selected as -1 (aborted)
		case "ctrl+c", "q", "esc":
			m.selected = -1
			return m, tea.Quit

		// Navigate cursor up (with wrapping support)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}

		// Navigate cursor down (with wrapping support)
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		// Confirm choice: save the current cursor index and exit the program loop
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

// View outputs the visual representation of the selection menu in the terminal.
// It uses package-private Lipgloss styles defined in wizard.go to preserve uniform visuals.
func (m menuModel) View() string {
	var s strings.Builder

	// Top title header panel
	s.WriteString(titleStyle.Render("   CONFIG MAKER - INTERACTIVE MENU   ") + "\n\n")

	// Render choices list
	var content strings.Builder
	content.WriteString(" Select Interface / Utility:\n\n")

	for i, choice := range m.choices {
		if i == m.cursor {
			// Styled pointer row for the currently selected/active item
			content.WriteString(fmt.Sprintf("  ▶ %s\n", activeItemStyle.Render(choice)))
		} else {
			// Dimmed representation for inactive choices
			content.WriteString(fmt.Sprintf("    %s\n", inactiveItemStyle.Render(choice)))
		}
	}

	// Bottom navigation help notes
	content.WriteString("\n" + helpStyle.Render("Use Up/Down (or j/k) to select, Enter to confirm, q to exit."))

	// Wrap entire content in a clean slate-gray box
	s.WriteString(boxStyle.Render(content.String()))
	return s.String()
}

// RunMenu launches the Bubble Tea loop for the main selection menu.
// It returns the index of the selected choice (0 to 4), or -1 if the user chose to exit.
func RunMenu() int {
	m := menuModel{
		choices: []string{
			"CLI Wizard (Bubble Tea TUI)",
			"Web Dashboard (Local Server)",
			"Import Settings (View config.json)",
			"Export Defaults (Reset config.json)",
			"Exit",
		},
		selected: -1,
	}

	// Start program execution loop
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return -1
	}

	// Extract selected choice index from final program model
	if fm, ok := finalModel.(menuModel); ok {
		return fm.selected
	}

	return -1
}
