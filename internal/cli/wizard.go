package cli

import (
	"config-maker/internal/config"
	"config-maker/internal/executor"
	"config-maker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss TUI Styles
var (
	cyan        = lipgloss.Color("15")  // Crisp white for active items and highlights
	purple      = lipgloss.Color("240") // Dark slate gray for subtle borders
	green       = lipgloss.Color("120")
	red         = lipgloss.Color("196")
	gray        = lipgloss.Color("243") // Muted gray for secondary items
	obsidianBg  = lipgloss.Color("233")

	titleStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple).
			Padding(1, 2).
			Width(74)

	activeItemStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	inactiveItemStyle = lipgloss.NewStyle().
			Foreground(gray)

	successStyle = lipgloss.NewStyle().
			Foreground(green)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	helpStyle = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true)
)

// Custom TUI Message definitions
type logLineMsg string
type execFinishedMsg struct {
	err error
}

type tuiModel struct {
	step             int // 0 to 7
	cfg              config.UserConfig
	hasSavedSettings bool
	importPrompt     bool // true = Yes, false = No

	// TUI selectors cursors
	ohMyZshChoice    bool // true = Yes, false = No
	configureGit     bool // true = Yes, false = No
	applyTheme       bool // true = Yes, false = No
	themeModeCursor  int  // 0 = Dark, 1 = Light
	themeCursor      int  // index in availableThemes
	availableThemes  []string
	applyBg          bool // true = Yes, false = No
	bgChoiceCursor   int  // 0 = Predefined, 1 = Repo List, 2 = Custom Path
	wallpaperCursor  int  // index in repoWallpapers
	repoWallpapers   []string
	enableDocker     bool // true = Yes, false = No
	enableZshDefault bool // true = Yes, false = No
	exportPrompt     bool // true = Yes, false = No
	importedSettings bool // true if settings were imported at Step 0


	// Form text inputs
	gitNameInput   textinput.Model
	gitEmailInput  textinput.Model
	customBgInput  textinput.Model
	focusedInput   int // 0 or 1 for Git, 0 for Custom Path

	// Execution logs state
	applying   bool
	finished   bool
	err        error
	logChan    chan tea.Msg
	logLines   []string
}

// RunWizard kicks off the interactive Bubble Tea terminal wizard.
func RunWizard() {
	homeDir, _ := os.UserHomeDir()
	configFilePath := filepath.Join(homeDir, ".config", "config-maker", "config.json")
	_, err := os.Stat(configFilePath)
	hasSaved := err == nil

	// Pre-load saved or default configuration
	var initCfg config.UserConfig
	if hasSaved {
		initCfg = config.LoadConfig()
	} else {
		initCfg = config.DefaultConfig()
	}

	// Initialize Git text inputs
	nameInput := textinput.New()
	nameInput.Placeholder = "e.g. John Doe"
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = 30
	if initCfg.GitName != "" {
		nameInput.SetValue(initCfg.GitName)
	}

	emailInput := textinput.New()
	emailInput.Placeholder = "e.g. john@example.com"
	emailInput.CharLimit = 64
	emailInput.Width = 30
	if initCfg.GitEmail != "" {
		emailInput.SetValue(initCfg.GitEmail)
	}

	// Initialize custom background path input
	customBgInput := textinput.New()
	customBgInput.Placeholder = "/home/user/Pictures/wallpaper.jpg"
	customBgInput.CharLimit = 256
	customBgInput.Width = 45
	if initCfg.BackgroundImage != "" && !strings.Contains(initCfg.BackgroundImage, "Zone01_Desk_cfg") {
		customBgInput.SetValue(initCfg.BackgroundImage)
	}

	// Initial scanning of themes and wallpapers
	themes := getSystemThemes("1")
	wallpapers := []string{
		"976013.jpg",
		"Rin_Shima_Level_Up_Your_Web_Apps_With_Go.png",
		"wallpaper-01.png",
	}

	m := tuiModel{
		step:             0,
		cfg:              initCfg,
		hasSavedSettings: hasSaved,
		importPrompt:     true,

		// Pre-populate selections from loaded config
		ohMyZshChoice:    initCfg.InstallOhMyZsh,
		configureGit:     initCfg.ConfigureGit,
		applyTheme:       initCfg.ApplyTheme,
		themeModeCursor:  0, // default dark
		themeCursor:      0,
		availableThemes:  themes,
		applyBg:          initCfg.ApplyBackground,
		bgChoiceCursor:   0,
		wallpaperCursor:  0,
		repoWallpapers:   wallpapers,
		enableDocker:     initCfg.EnableDocker,
		enableZshDefault: initCfg.EnableZshDefault,
		exportPrompt:     true,

		gitNameInput:  nameInput,
		gitEmailInput: emailInput,
		customBgInput: customBgInput,
		focusedInput:  0,
		logChan:       make(chan tea.Msg, 100),
	}

	if initCfg.ThemeMode == "2" {
		m.themeModeCursor = 1
	}

	// Resolve the background image type based on loaded path
	if initCfg.BackgroundImage != "" {
		if strings.Contains(initCfg.BackgroundImage, "wallpapers/") {
			m.bgChoiceCursor = 1
			filename := filepath.Base(initCfg.BackgroundImage)
			for idx, wp := range wallpapers {
				if wp == filename {
					m.wallpaperCursor = idx
					break
				}
			}
		} else if !strings.Contains(initCfg.BackgroundImage, "Background.jpeg") {
			m.bgChoiceCursor = 2
		}
	}

	// If the user has no saved settings, skip Step 0 (Import saved settings prompt)
	if !hasSaved {
		m.step = 1
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf(" TUI error: %v\n", err)
		os.Exit(1)
	}

	// Trigger terminal restart/reload after TUI finishes successfully
	if fm, ok := finalModel.(tuiModel); ok && fm.finished && fm.err == nil {
		executor.FinishSetup(os.Stdout)
	}
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

// waitForActivity command reads asynchronously from the logging channel.
func waitForActivity(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// executeApplyConfig is a tea.Cmd wrapper that runs ApplyConfig in a goroutine.
func executeApplyConfig(cfg config.UserConfig, export bool, sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		writer := &tuiLogWriter{sub: sub}
		err := executor.ApplyConfig(cfg, export, writer)
		return execFinishedMsg{err: err}
	}
}

type tuiLogWriter struct {
	sub chan tea.Msg
}

func (tlw *tuiLogWriter) Write(p []byte) (int, error) {
	lines := strings.Split(string(p), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) > 0 {
			// Strip raw ANSI color escape codes for clean terminal TUI lines
			trimmed = stripAnsiColors(trimmed)
			tlw.sub <- logLineMsg(trimmed)
		}
	}
	return len(p), nil
}

func stripAnsiColors(s string) string {
	r := strings.NewReplacer(
		"\033[0m", "", "\033[0;31m", "", "\033[0;32m", "",
		"\033[0;33m", "", "\033[0;36m", "", "\033[1;0;31m", "",
		"\033[1;0;32m", "", "\033[1;0;33m", "", "\033[1;0;36m", "",
		"\033[1m\033[92m", "",
	)
	return r.Replace(s)
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global abort
		if msg.String() == "ctrl+c" || (m.step != 7 && (msg.String() == "q" || msg.String() == "esc")) {
			if m.applying && !m.finished {
				return m, nil // prevent quit during execution
			}
			return m, tea.Quit
		}

		if m.applying {
			if m.finished {
				if msg.String() == "enter" || msg.String() == "q" || msg.String() == "esc" {
					return m, tea.Quit
				}
			}
			return m, nil // lock keys during apply
		}

		// Handle key presses based on current step
		switch m.step {
		case 0: // Import Saved settings prompt
			switch msg.String() {
			case "left", "right", "tab", "h", "l", "up", "down", "j", "k":
				m.importPrompt = !m.importPrompt
			case "enter":
				if m.importPrompt {
					m.cfg = config.LoadConfig()
					m.ohMyZshChoice = m.cfg.InstallOhMyZsh
					m.configureGit = m.cfg.ConfigureGit
					m.applyTheme = m.cfg.ApplyTheme
					m.enableDocker = m.cfg.EnableDocker
					m.enableZshDefault = m.cfg.EnableZshDefault
					m.gitNameInput.SetValue(m.cfg.GitName)
					m.gitEmailInput.SetValue(m.cfg.GitEmail)
					if m.cfg.ThemeMode == "2" {
						m.themeModeCursor = 1
					} else {
						m.themeModeCursor = 0
					}
					m.availableThemes = getSystemThemes(m.cfg.ThemeMode)
					// resolve background selection
					if m.cfg.BackgroundImage != "" {
						if strings.Contains(m.cfg.BackgroundImage, "wallpapers/") {
							m.bgChoiceCursor = 1
							m.wallpaperCursor = 0
							wpFilename := filepath.Base(m.cfg.BackgroundImage)
							for idx, wp := range m.repoWallpapers {
								if wp == wpFilename {
									m.wallpaperCursor = idx
									break
								}
							}
						} else if !strings.Contains(m.cfg.BackgroundImage, "Background.jpeg") {
							m.bgChoiceCursor = 2
							m.customBgInput.SetValue(m.cfg.BackgroundImage)
						} else {
							m.bgChoiceCursor = 0
						}
					}
					m.importedSettings = true
					m.step = 4 // Skip straight to background selection!
				} else {
					m.cfg = config.DefaultConfig()
					m.importedSettings = false
					m.step = 1
				}
			}

		case 1: // Oh-My-Zsh Choice
			switch msg.String() {
			case "backspace":
				if m.hasSavedSettings {
					m.step = 0
				}
			case "left", "right", "tab", "h", "l", "up", "down", "j", "k":
				m.ohMyZshChoice = !m.ohMyZshChoice
			case "enter":
				m.step = 2
			}

		case 2: // Git Config Setup
			if !m.configureGit {
				switch msg.String() {
				case "backspace":
					m.step = 1
				case "left", "right", "tab", "h", "l", "up", "down", "j", "k":
					m.configureGit = !m.configureGit
				case "enter":
					m.step = 3
				}
			} else {
				// Text input focus and typing
				var cmd tea.Cmd
				if m.focusedInput == 0 {
					switch msg.String() {
					case "down", "tab":
						m.focusedInput = 1
						m.gitNameInput.Blur()
						m.gitEmailInput.Focus()
					case "enter":
						m.focusedInput = 1
						m.gitNameInput.Blur()
						m.gitEmailInput.Focus()
					default:
						m.gitNameInput, cmd = m.gitNameInput.Update(msg)
					}
				} else if m.focusedInput == 1 {
					switch msg.String() {
					case "up", "tab":
						m.focusedInput = 0
						m.gitEmailInput.Blur()
						m.gitNameInput.Focus()
					case "enter":
						// save credentials and move next
						m.cfg.GitName = m.gitNameInput.Value()
						m.cfg.GitEmail = m.gitEmailInput.Value()
						m.step = 3
					default:
						m.gitEmailInput, cmd = m.gitEmailInput.Update(msg)
					}
				}
				return m, cmd
			}

		case 3: // Themes Setup
			if !m.applyTheme {
				switch msg.String() {
				case "backspace":
					m.step = 2
				case "left", "right", "tab", "h", "l", "up", "down", "j", "k":
					m.applyTheme = !m.applyTheme
				case "enter":
					m.step = 4
				}
			} else {
				switch msg.String() {
				case "backspace":
					m.step = 2
				case "left", "right", "tab", "h", "l":
					m.themeModeCursor = 1 - m.themeModeCursor
					modeStr := "1"
					if m.themeModeCursor == 1 {
						modeStr = "2"
					}
					m.availableThemes = getSystemThemes(modeStr)
					m.themeCursor = 0
				case "up", "k":
					if m.themeCursor > 0 {
						m.themeCursor--
					}
				case "down", "j":
					if m.themeCursor < len(m.availableThemes)-1 {
						m.themeCursor++
					}
				case "enter":
					if len(m.availableThemes) > 0 {
						m.cfg.ThemeName = m.availableThemes[m.themeCursor]
					} else {
						m.cfg.ThemeName = "Yaru-dark"
					}
					if m.themeModeCursor == 1 {
						m.cfg.ThemeMode = "2"
					} else {
						m.cfg.ThemeMode = "1"
					}
					m.step = 4
				}
			}

		case 4: // Background Setup
			if !m.applyBg {
				switch msg.String() {
				case "backspace":
					if m.importedSettings {
						m.step = 0
					} else {
						m.step = 3
					}
				case "left", "right", "tab", "h", "l", "up", "down", "j", "k":
					m.applyBg = !m.applyBg
				case "enter":
					if m.importedSettings {
						m.step = 6
					} else {
						m.step = 5
					}
				}
			} else {
				// We have a background choice selection
				if m.bgChoiceCursor == 2 {
					// Custom path typing
					var cmd tea.Cmd
					m.customBgInput.Focus()
					if msg.String() == "enter" {
						m.cfg.BackgroundImage = m.customBgInput.Value()
						if m.importedSettings {
							m.step = 6
						} else {
							m.step = 5
						}
					} else if msg.String() == "up" {
						m.bgChoiceCursor = 1
						m.customBgInput.Blur()
					} else {
						m.customBgInput, cmd = m.customBgInput.Update(msg)
					}
					return m, cmd
				} else {
					switch msg.String() {
					case "backspace":
						if m.importedSettings {
							m.step = 0
						} else {
							m.step = 3
						}
					case "up", "k":
						if m.bgChoiceCursor > 0 {
							m.bgChoiceCursor--
						}
					case "down", "j":
						if m.bgChoiceCursor < 2 {
							m.bgChoiceCursor++
						}
					case "left", "h":
						if m.bgChoiceCursor == 1 && m.wallpaperCursor > 0 {
							m.wallpaperCursor--
						}
					case "right", "l":
						if m.bgChoiceCursor == 1 && m.wallpaperCursor < len(m.repoWallpapers)-1 {
							m.wallpaperCursor++
						}
					case "enter":
						homeDir, _ := os.UserHomeDir()
						if m.bgChoiceCursor == 0 {
							m.cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "Background.jpeg")
						} else if m.bgChoiceCursor == 1 {
							selectedWP := m.repoWallpapers[m.wallpaperCursor]
							m.cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "wallpapers", selectedWP)
						}
						if m.importedSettings {
							m.step = 6
						} else {
							m.step = 5
						}
					}
				}
			}

		case 5: // Docker and Zsh default selectors
			// We have two toggle switches: active switch cursor (focusedInput 0 or 1)
			switch msg.String() {
			case "backspace":
				m.step = 4
			case "up", "down", "tab":
				m.focusedInput = 1 - m.focusedInput
			case "left", "right", "h", "l", " ":
				if m.focusedInput == 0 {
					m.enableDocker = !m.enableDocker
				} else {
					m.enableZshDefault = !m.enableZshDefault
				}
			case "enter":
				m.cfg.InstallOhMyZsh = m.ohMyZshChoice
				m.cfg.ConfigureGit = m.configureGit
				m.cfg.ApplyTheme = m.applyTheme
				m.cfg.ApplyBackground = m.applyBg
				m.cfg.EnableDocker = m.enableDocker
				m.cfg.EnableZshDefault = m.enableZshDefault
				m.step = 6
			}

		case 6: // Final Summary & Apply
			switch msg.String() {
			case "backspace":
				if m.importedSettings {
					m.step = 4
				} else {
					m.step = 5
				}
			case "tab", "left", "right", "h", "l", "up", "down", "j", "k":
				m.exportPrompt = !m.exportPrompt
			case "enter":
				m.applying = true
				m.step = 7
				// Launch native Go ApplyConfig in a concurrent goroutine!
				return m, tea.Batch(
					executeApplyConfig(m.cfg, m.exportPrompt, m.logChan),
					waitForActivity(m.logChan),
				)
			}
		}

	// Logging console channels updates
	case logLineMsg:
		m.logLines = append(m.logLines, string(msg))
		return m, waitForActivity(m.logChan)

	case execFinishedMsg:
		m.finished = true
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

func (m tuiModel) View() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("   CONFIG MAKER - TERMINAL UI WIZARD   ") + "\n\n")

	var stepContent string

	switch m.step {
	case 0: // Import Settings Prompt
		stepContent = fmt.Sprintf(
			" Loaded existing configurations in ~/.config/config-maker/config.json.\n\n"+
				" %sWould you like to import your saved settings?%s\n\n"+
				"   %s\n   %s",
			utils.Yellow, utils.Reset,
			renderToggleOption("Yes (Import Settings)", m.importPrompt),
			renderToggleOption("No (Pristine Defaults)", !m.importPrompt),
		)

	case 1: // Oh-My-Zsh Setup
		stepContent = fmt.Sprintf(
			" [Step 1/6] Oh-My-Zsh Installation\n\n"+
				" %sWould you like to install Oh-My-Zsh unattended?%s\n\n"+
				"   %s\n   %s",
			utils.Yellow, utils.Reset,
			renderToggleOption("Yes, install Oh-My-Zsh", m.ohMyZshChoice),
			renderToggleOption("No, skip", !m.ohMyZshChoice),
		)

	case 2: // Git Setup
		if !m.configureGit {
			stepContent = fmt.Sprintf(
				" [Step 2/6] Git global setup\n\n"+
					" %sWould you like to configure Git name & email?%s\n\n"+
					"   %s\n   %s",
					utils.Yellow, utils.Reset,
				renderToggleOption("Yes, configure Git", m.configureGit),
				renderToggleOption("No, skip Git setup", !m.configureGit),
			)
		} else {
			stepContent = fmt.Sprintf(
				" [Step 2/6] Enter Git credentials\n\n"+
					"   %s Name  %s\n"+
					"   %s\n\n"+
					"   %s Email %s\n"+
					"   %s\n\n"+
					"   %s",
				renderActiveLabel("Full Name", m.focusedInput == 0), utils.Reset,
				m.gitNameInput.View(),
				renderActiveLabel("Email Address", m.focusedInput == 1), utils.Reset,
				m.gitEmailInput.View(),
				helpStyle.Render("Press Down/Tab to switch fields, and Enter to complete."),
			)
		}

	case 3: // Themes Setup
		if !m.applyTheme {
			stepContent = fmt.Sprintf(
				" [Step 3/6] Gnome visual theme\n\n"+
					" %sWould you like to customize Gnome windows/interface theme?%s\n\n"+
					"   %s\n   %s",
					utils.Yellow, utils.Reset,
				renderToggleOption("Yes, apply Gnome themes", m.applyTheme),
				renderToggleOption("No, skip themes", !m.applyTheme),
			)
		} else {
			// Mode selection and theme listing
			modeView := renderHorizontalOptions([]string{"Dark Mode", "Light Mode"}, m.themeModeCursor)

			var themesView strings.Builder
			themesView.WriteString(fmt.Sprintf("   Select System Theme:\n"))
			if len(m.availableThemes) == 0 {
				themesView.WriteString(fmt.Sprintf("     %sNo themes found on system.%s\n", utils.Red, utils.Reset))
			} else {
				start := m.themeCursor - 2
				if start < 0 {
					start = 0
				}
				end := start + 5
				if end > len(m.availableThemes) {
					end = len(m.availableThemes)
				}
				for i := start; i < end; i++ {
					if i == m.themeCursor {
						themesView.WriteString(fmt.Sprintf("     %s▶ %s%s\n", utils.Cyan, m.availableThemes[i], utils.Reset))
					} else {
						themesView.WriteString(fmt.Sprintf("       %s\n", m.availableThemes[i]))
					}
				}
			}

			stepContent = fmt.Sprintf(
				" [Step 3/6] Select Gnome style\n\n"+
					"   Select Mode: %s\n\n"+
					"%s\n"+
					"   %s",
				modeView,
				themesView.String(),
				helpStyle.Render("Press Left/Right for mode, Up/Down for themes, Enter to save."),
			)
		}

	case 4: // Background Setup
		if !m.applyBg {
			stepContent = fmt.Sprintf(
				" [Step 4/6] Desktop Wallpaper\n\n"+
					" %sWould you like to set custom desktop backgrounds?%s\n\n"+
					"   %s\n   %s",
					utils.Yellow, utils.Reset,
				renderToggleOption("Yes, apply desktop wallpaper", m.applyBg),
				renderToggleOption("No, skip wallpaper", !m.applyBg),
			)
		} else {
			// Render wallpaper option selections
			var list string
			list += renderOptionRow("[1] Predefined Background.jpeg", m.bgChoiceCursor == 0) + "\n"

			// Choice 2: wallpapers scrollable list inline
			wpSelector := ""
			if m.bgChoiceCursor == 1 {
				wpSelector = "  ◀ " + m.repoWallpapers[m.wallpaperCursor] + " ▶"
			}
			list += renderOptionRow("[2] Choice from repository wallpapers"+wpSelector, m.bgChoiceCursor == 1) + "\n"

			list += renderOptionRow("[3] Custom absolute image path", m.bgChoiceCursor == 2) + "\n"

			var pathInputView string
			if m.bgChoiceCursor == 2 {
				pathInputView = fmt.Sprintf("\n   Enter Path: %s\n", m.customBgInput.View())
			}

			stepContent = fmt.Sprintf(
				" [Step 4/6] Wallpaper choice\n\n"+
					"%s"+
					"%s\n"+
					"   %s",
				list,
				pathInputView,
				helpStyle.Render("Press Up/Down to choose option, Left/Right to scroll wallpapers, Enter to select."),
			)
		}

	case 5: // Docker and Zsh defaults
		stepContent = fmt.Sprintf(
			" [Step 5/6] Environments & defaults\n\n"+
				"   %s Enable Docker Rootless?%s\n"+
				"     %s\n\n"+
				"   %s Set Zsh as Default Shell?%s\n"+
				"     %s\n\n"+
				"   %s",
			renderActiveLabel("[Option A]", m.focusedInput == 0), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.enableDocker)),
			renderActiveLabel("[Option B]", m.focusedInput == 1), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.enableZshDefault)),
			helpStyle.Render("Press Up/Down/Tab to switch, Left/Right/Space to toggle, Enter to confirm."),
		)

	case 6: // Final Summary Checklist & Save
		var summary strings.Builder
		summary.WriteString(fmt.Sprintf(" %sReview Selections:%s\n\n", utils.Yellow, utils.Reset))
		summary.WriteString(renderSummaryRow("Install Oh-My-Zsh", m.ohMyZshChoice) + "\n")
		summary.WriteString(renderSummaryRow("Configure Git", m.configureGit) + "\n")
		if m.configureGit {
			summary.WriteString(fmt.Sprintf("   └─ Name:  %s\n", m.gitNameInput.Value()))
			summary.WriteString(fmt.Sprintf("   └─ Email: %s\n", m.gitEmailInput.Value()))
		}
		summary.WriteString(renderSummaryRow("Apply theme", m.applyTheme) + "\n")
		if m.applyTheme {
			mode := "Dark"
			if m.themeModeCursor == 1 {
				mode = "Light"
			}
			theme := "Yaru-dark"
			if len(m.availableThemes) > 0 {
				theme = m.availableThemes[m.themeCursor]
			}
			summary.WriteString(fmt.Sprintf("   └─ %s (%s)\n", theme, mode))
		}
		summary.WriteString(renderSummaryRow("Apply Background", m.applyBg) + "\n")
		summary.WriteString(renderSummaryRow("Install Docker Rootless", m.enableDocker) + "\n")
		summary.WriteString(renderSummaryRow("Set Zsh Default Shell", m.enableZshDefault) + "\n\n")

		summary.WriteString(fmt.Sprintf(" %sWould you like to export/save these settings for future use?%s\n", utils.Yellow, utils.Reset))
		summary.WriteString(fmt.Sprintf("   %s\n   %s\n\n",
			renderToggleOption("Yes, save choices", m.exportPrompt),
			renderToggleOption("No, do not save", !m.exportPrompt),
		))

		summary.WriteString(helpStyle.Render("Press Enter to execute the configuration updates, or Esc to quit."))
		stepContent = summary.String()

	case 7: // Live Execution Console
		var logs strings.Builder
		logs.WriteString(fmt.Sprintf(" %sExecution Console:%s\n\n", utils.Yellow, utils.Reset))

		if len(m.logLines) == 0 {
			logs.WriteString("   Initialising setup environment...\n")
		} else {
			start := len(m.logLines) - 10
			if start < 0 {
				start = 0
			}
			for i := start; i < len(m.logLines); i++ {
				logs.WriteString("   " + m.logLines[i] + "\n")
			}
		}

		logs.WriteString("\n")

		if !m.finished {
			logs.WriteString(fmt.Sprintf("   %s● Applying configuration selections live... Please wait.%s\n", utils.Cyan, utils.Reset))
		} else {
			if m.err != nil {
				logs.WriteString(fmt.Sprintf("   %s✖ Failed: %v%s\n\n", utils.Red, m.err, utils.Reset))
				logs.WriteString("   Press [q] or [Esc] to exit.")
			} else {
				logs.WriteString(fmt.Sprintf("   %s✔ Finished successfully!%s\n\n", utils.Green, utils.Reset))
				logs.WriteString("   Press [Enter] to reload and reopen GNOME Terminals now.")
			}
		}
		stepContent = logs.String()
	}

	// If we are finished successfully in step 7 and press Enter, exit & reload Gnome terminal
	if m.finished && m.err == nil && m.step == 7 {
		// Finish behavior will trigger on tea.Quit, let's catch it in the caller or exit.
		// Wait! Since we are in the View, we should return immediately on tea.Quit.
	}

	s.WriteString(boxStyle.Render(stepContent))
	return s.String()
}

// Render UI element helpers
func renderToggleOption(label string, active bool) string {
	if active {
		return fmt.Sprintf("%s● %s%s", utils.Cyan, label, utils.Reset)
	}
	return fmt.Sprintf("  %s", label)
}

func renderOptionRow(label string, active bool) string {
	if active {
		return fmt.Sprintf("  %s▶ %s%s", utils.Cyan, label, utils.Reset)
	}
	return fmt.Sprintf("    %s", label)
}

func renderHorizontalOptions(options []string, activeIndex int) string {
	var sb strings.Builder
	for i, opt := range options {
		if i == activeIndex {
			sb.WriteString(fmt.Sprintf("%s(•) %s%s   ", utils.Cyan, opt, utils.Reset))
		} else {
			sb.WriteString(fmt.Sprintf("( ) %s   ", opt))
		}
	}
	return sb.String()
}

func renderSummaryRow(label string, enabled bool) string {
	if enabled {
		return fmt.Sprintf("   %s✔%s %s", utils.Green, utils.Reset, label)
	}
	return fmt.Sprintf("   %s✖%s %s", utils.Red, utils.Reset, label)
}

func getYesNoIndex(no bool) int {
	if no {
		return 1
	}
	return 0
}

func renderActiveLabel(label string, active bool) string {
	if active {
		return fmt.Sprintf("%s%s%s", utils.Cyan, label, utils.Reset)
	}
	return label
}

// getSystemThemes helper returns matching themes
func getSystemThemes(mode string) []string {
	var themes []string
	entries, err := os.ReadDir("/usr/share/themes")
	if err != nil {
		return []string{"Yaru-dark"}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if mode == "1" && strings.Contains(strings.ToLower(name), "dark") {
			themes = append(themes, name)
		} else if mode == "2" && !strings.Contains(strings.ToLower(name), "dark") {
			if name == "Default" || name == "raleigh" {
				continue
			}
			themes = append(themes, name)
		}
	}
	return themes
}

// getRepositoryWallpapers helper returns standard values
func getRepositoryWallpapers() []string {
	return []string{
		"976013.jpg",
		"Rin_Shima_Level_Up_Your_Web_Apps_With_Go.png",
		"wallpaper-01.png",
	}
}
