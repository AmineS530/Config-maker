package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zonerestore/internal/config"
	"zonerestore/internal/executor"
	"zonerestore/internal/themes"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RunWizard kicks off the interactive Bubble Tea terminal wizard.
func RunWizard() {
	homeDir, _ := os.UserHomeDir()
	configFilePath := filepath.Join(homeDir, ".config", "zonerestore", "config.json")
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
	nameInput.Placeholder = "e.g. username"
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = 30
	if initCfg.Git.GitName != "" {
		nameInput.SetValue(initCfg.Git.GitName)
	}

	emailInput := textinput.New()
	emailInput.Placeholder = "e.g. username@example.com"
	emailInput.CharLimit = 64
	emailInput.Width = 30
	if initCfg.Git.GitEmail != "" {
		emailInput.SetValue(initCfg.Git.GitEmail)
	}

	// Initialize custom background path input
	customBgInput := textinput.New()
	customBgInput.Placeholder = "/home/user/Pictures/wallpaper.jpg"
	customBgInput.CharLimit = 256
	customBgInput.Width = 45
	if initCfg.Wallpaper.BackgroundImage != "" && !strings.Contains(initCfg.Wallpaper.BackgroundImage, "ZoneRestoreThemes") {
		customBgInput.SetValue(initCfg.Wallpaper.BackgroundImage)
	}

	// Initialize themes and wallpapers dynamically from repository
	themesRoot := themes.Root()
	repoWallpapers := themes.ListWallpapers(themesRoot)
	if len(repoWallpapers) == 0 {
		repoWallpapers = []string{
			"976013.jpg",
			"Rin_Shima_Level_Up_Your_Web_Apps_With_Go.png",
			"wallpaper-01.png",
		}
	}

	repoFonts := themes.ListFonts(themesRoot)
	if len(repoFonts) == 0 {
		repoFonts = []string{
			"MesloLGS NF",
			"Fira Code",
			"JetBrains Mono",
		}
	}

	systemThemes := getSystemThemes("1")

	// Set up execution console spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(brandCyan)

	m := tuiModel{
		step:             0,
		cfg:              initCfg,
		hasSavedSettings: hasSaved,
		importPrompt:     true,

		// Pre-populate selections from loaded config
		ohMyZshChoice:   initCfg.Zsh.InstallOhMyZsh,
		configureGit:    initCfg.Git.ConfigureGit,
		applyTheme:      initCfg.Theme.ApplyTheme,
		themeModeCursor: 0, // default dark
		themeCursor:      0,
		availableThemes:  systemThemes,

		applyFonts:     initCfg.Fonts.ConfigureFonts,
		fontCursor:     0,
		availableFonts: repoFonts,

		applyBg:         initCfg.Wallpaper.ApplyBackground,
		bgChoiceCursor:  0,
		wallpaperCursor: 0,
		repoWallpapers:  repoWallpapers,

		enableDocker:      initCfg.Docker.EnableDocker,
		enableZshDefault:  initCfg.Shell.EnableZshDefault,
		pinDiscord:        initCfg.Dock.PinDiscord,
		configureKeyboard: initCfg.Keyboard.ConfigureKeyboard,
		addArabic:         initCfg.Keyboard.AddArabic,

		exportPrompt: true,

		gitNameInput:  nameInput,
		gitEmailInput: emailInput,
		customBgInput: customBgInput,
		focusedInput:  0,
		logChan:       make(chan tea.Msg, 100),
		spinner:       sp,
	}

	if initCfg.Theme.ThemeMode == "2" {
		m.themeModeCursor = 1
	}

	// Resolve the background image type based on loaded path
	if initCfg.Wallpaper.BackgroundImage != "" {
		if strings.Contains(initCfg.Wallpaper.BackgroundImage, "themes/wallpapers/") {
			m.bgChoiceCursor = 1
			filename := filepath.Base(initCfg.Wallpaper.BackgroundImage)
			for idx, wp := range repoWallpapers {
				if wp == filename {
					m.wallpaperCursor = idx
					break
				}
			}
		} else if !strings.Contains(initCfg.Wallpaper.BackgroundImage, "Background.jpeg") {
			m.bgChoiceCursor = 2
		}
	}

	// Resolve the display font selection cursor from saved config
	// FontName (terminal) is always MesloLGS NF — never changed by the picker
	m.cfg.Fonts.FontName = "MesloLGS NF"
	if initCfg.Fonts.DisplayFontName != "" {
		for idx, f := range repoFonts {
			if f == initCfg.Fonts.DisplayFontName {
				m.fontCursor = idx + 1 // +1 because cursor 0 = "System Default"
				break
			}
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
		if trimmed != "" {
			// Strip ANSI escape colors to keep local logs clean
			trimmed = strings.ReplaceAll(trimmed, "\r", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[0m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[0;31m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[0;32m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[0;33m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[0;36m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[1;0;31m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[1;0;32m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[1;0;33m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[1;0;36m", "")
			trimmed = strings.ReplaceAll(trimmed, "\033[1m\033[92m", "")
			tlw.sub <- logLineMsg(trimmed)
		}
	}
	return len(p), nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Master exit keys (anytime unless running configuration)
		if (msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "esc") && !m.applying {
			return m, tea.Quit
		}

		// When applying, spinner continues running
		if m.applying && !m.finished {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

		// Save/Export step after finished
		if m.finished && m.err == nil {
			switch msg.String() {
			case "up", "down", "tab", "j", "k":
				m.exportPrompt = !m.exportPrompt
			case "enter":
				m.exportDone = true
				if m.exportPrompt {
					// Re-run apply with export settings true (just saves file instantly)
					_ = config.SaveConfig(m.cfg)
				}
				return m, tea.Quit
			}
			return m, nil
		}

		// Regular Wizard steps navigation
		switch m.step {
		case 0: // Saved settings import
			switch msg.String() {
			case "up", "down", "j", "k", "tab":
				m.importPrompt = !m.importPrompt
			case "enter":
				if m.importPrompt {
					m.cfg = config.LoadConfig()
					m.importedSettings = true
					// Prepopulate state from imported config
					m.ohMyZshChoice = m.cfg.Zsh.InstallOhMyZsh
					m.configureGit = m.cfg.Git.ConfigureGit
					m.applyTheme = m.cfg.Theme.ApplyTheme
					m.applyFonts = m.cfg.Fonts.ConfigureFonts
					m.applyBg = m.cfg.Wallpaper.ApplyBackground
					m.enableDocker = m.cfg.Docker.EnableDocker
					m.enableZshDefault = m.cfg.Shell.EnableZshDefault
					m.pinDiscord = m.cfg.Dock.PinDiscord
					m.configureKeyboard = m.cfg.Keyboard.ConfigureKeyboard
					m.addArabic = m.cfg.Keyboard.AddArabic
					m.gitNameInput.SetValue(m.cfg.Git.GitName)
					m.gitEmailInput.SetValue(m.cfg.Git.GitEmail)
					if m.cfg.Wallpaper.BackgroundImage != "" {
						m.customBgInput.SetValue(m.cfg.Wallpaper.BackgroundImage)
					}
					// Instantly jump to final summary checklist screen
					m.step = 7
				} else {
					m.step = 1
				}
			}

		case 1: // Oh-My-Zsh choice
			switch msg.String() {
			case "backspace":
				if m.hasSavedSettings {
					m.step = 0
				}
			case "up", "down", "j", "k", "tab":
				m.ohMyZshChoice = !m.ohMyZshChoice
			case "enter":
				m.cfg.Zsh.InstallOhMyZsh = m.ohMyZshChoice
				m.step = 2
			}

		case 2: // Git credentials
			if !m.configureGit {
				switch msg.String() {
				case "backspace":
					m.step = 1
				case "up", "down", "j", "k", "tab":
					m.configureGit = !m.configureGit
				case "enter":
					m.cfg.Git.ConfigureGit = false
					m.step = 3
				}
			} else {
				// Form inputs are active
				var cmd tea.Cmd
				if m.focusedInput == 0 {
					m.gitNameInput.Focus()
					m.gitEmailInput.Blur()
				} else {
					m.gitNameInput.Blur()
					m.gitEmailInput.Focus()
				}

				switch msg.String() {
				case "backspace":
					if m.gitNameInput.Focused() && m.gitNameInput.Value() == "" {
						m.configureGit = false
					} else if m.gitEmailInput.Focused() && m.gitEmailInput.Value() == "" {
						m.focusedInput = 0
					} else {
						if m.focusedInput == 0 {
							m.gitNameInput, cmd = m.gitNameInput.Update(msg)
						} else {
							m.gitEmailInput, cmd = m.gitEmailInput.Update(msg)
						}
					}
				case "up", "k":
					m.focusedInput = 0
				case "down", "j", "tab":
					if m.focusedInput == 0 {
						m.focusedInput = 1
					} else {
						m.focusedInput = 0
					}
				case "enter":
					if m.focusedInput == 0 {
						m.focusedInput = 1
					} else {
						m.cfg.Git.ConfigureGit = true
						m.cfg.Git.GitName = m.gitNameInput.Value()
						m.cfg.Git.GitEmail = m.gitEmailInput.Value()
						m.step = 3
					}
				default:
					if m.focusedInput == 0 {
						m.gitNameInput, cmd = m.gitNameInput.Update(msg)
					} else {
						m.gitEmailInput, cmd = m.gitEmailInput.Update(msg)
					}
				}
				return m, cmd
			}

		case 3: // Themes choice
			if !m.applyTheme {
				switch msg.String() {
				case "backspace":
					m.step = 2
				case "up", "down", "j", "k", "tab":
					m.applyTheme = !m.applyTheme
				case "enter":
					m.cfg.Theme.ApplyTheme = false
					m.step = 4
				}
			} else {
				// Theme settings active
				switch msg.String() {
				case "backspace":
					m.applyTheme = false
				case "left", "right", "h", "l":
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
					m.cfg.Theme.ApplyTheme = true
					if m.themeModeCursor == 0 {
						m.cfg.Theme.ThemeMode = "1"
					} else {
						m.cfg.Theme.ThemeMode = "2"
					}
					if len(m.availableThemes) > 0 {
						m.cfg.Theme.ThemeName = m.availableThemes[m.themeCursor]
					}
					m.step = 4
				}
			}

		case 4: // Font Selection Setup
			if !m.applyFonts {
				switch msg.String() {
				case "backspace":
					m.step = 3
				case "up", "down", "j", "k", "tab":
					m.applyFonts = !m.applyFonts
				case "enter":
					m.cfg.Fonts.ConfigureFonts = false
					m.step = 5
				}
			} else {
				switch msg.String() {
				case "backspace":
					m.step = 3
				case "up", "k":
					if m.fontCursor > 0 {
						m.fontCursor--
					}
				case "down", "j":
					if m.fontCursor < len(m.availableFonts)-1 {
						m.fontCursor++
					}
				case "enter":
					// Terminal font is always MesloLGS NF — locked
					m.cfg.Fonts.FontName = "MesloLGS NF"
					// cursor 0 means "System Default", cursor > 0 means a repo font
					if m.fontCursor > 0 && len(m.availableFonts) >= m.fontCursor {
						m.cfg.Fonts.DisplayFontName = m.availableFonts[m.fontCursor-1]
					} else {
						m.cfg.Fonts.DisplayFontName = "" // empty = system default (Ubuntu 11)
					}
					m.cfg.Fonts.ConfigureFonts = true
					m.step = 5
				}
			}

		case 5: // Background Setup
			if !m.applyBg {
				switch msg.String() {
				case "backspace":
					m.step = 4
				case "left", "right", "tab", "h", "l", "up", "down", "j", "k":
					m.applyBg = !m.applyBg
				case "enter":
					m.cfg.Wallpaper.ApplyBackground = false
					if m.importedSettings {
						m.step = 7
					} else {
						m.step = 6
					}
				}
			} else {
				// We have a background choice selection
				if m.bgChoiceCursor == 2 {
					// Custom path typing
					var cmd tea.Cmd
					m.customBgInput.Focus()
					if msg.String() == "enter" {
						m.cfg.Wallpaper.BackgroundImage = m.customBgInput.Value()
						m.cfg.Wallpaper.ApplyBackground = true
						if m.importedSettings {
							m.step = 7
						} else {
							m.step = 6
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
						m.step = 4
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
						if m.bgChoiceCursor == 0 {
							m.cfg.Wallpaper.BackgroundImage = "Background.jpeg"
						} else if m.bgChoiceCursor == 1 {
							m.cfg.Wallpaper.BackgroundImage = m.repoWallpapers[m.wallpaperCursor]
						}
						m.cfg.Wallpaper.ApplyBackground = true
						if m.importedSettings {
							m.step = 7
						} else {
							m.step = 6
						}
					}
				}
			}

		case 6: // Modular checklists (Docker, Zsh default, Keyboard, Dock Favorites, Arabic, Power)
			switch msg.String() {
			case "backspace":
				m.step = 5
			case "up", "k":
				if m.focusedInput > 0 {
					m.focusedInput--
				} else {
					m.focusedInput = 5
				}
			case "down", "j", "tab":
				if m.focusedInput < 5 {
					m.focusedInput++
				} else {
					m.focusedInput = 0
				}
			case "left", "right", "h", "l", " ":
				switch m.focusedInput {
				case 0:
					m.enableDocker = !m.enableDocker
				case 1:
					m.enableZshDefault = !m.enableZshDefault
				case 2:
					m.pinDiscord = !m.pinDiscord
				case 3:
					m.configureKeyboard = !m.configureKeyboard
				case 4:
					m.addArabic = !m.addArabic
				}
			case "enter":
				m.cfg.Zsh.InstallOhMyZsh = m.ohMyZshChoice
				m.cfg.Git.ConfigureGit = m.configureGit
				m.cfg.Theme.ApplyTheme = m.applyTheme
				m.cfg.Fonts.ConfigureFonts = m.applyFonts
				m.cfg.Wallpaper.ApplyBackground = m.applyBg
				m.cfg.Docker.EnableDocker = m.enableDocker
				m.cfg.Shell.EnableZshDefault = m.enableZshDefault
				m.cfg.Dock.PinDiscord = m.pinDiscord
				m.cfg.Keyboard.ConfigureKeyboard = m.configureKeyboard
				m.cfg.Keyboard.AddArabic = m.addArabic
				m.step = 7
			}

		case 7: // Final Summary & Apply
			switch msg.String() {
			case "backspace":
				if m.importedSettings {
					m.step = 5
				} else {
					m.step = 6
				}
			case "enter":
				m.applying = true
				m.step = 8
				// Launch native Go ApplyConfig in a concurrent goroutine!
				return m, tea.Batch(
					m.spinner.Tick,
					executeApplyConfig(m.cfg, false, m.logChan),
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
