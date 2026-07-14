package cli

import (
	"zonerestore/internal/config"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss TUI Brand Styles
var (
	brandCyan    = lipgloss.Color("#00D4FF") // Vibrant cyan for highlight/selected
	brandPurple  = lipgloss.Color("#9B59B6") // Rich purple for borders
	brandGreen   = lipgloss.Color("#2ECC71") // Success green
	brandRed     = lipgloss.Color("#E74C3C") // Error red
	brandGray    = lipgloss.Color("#7F8C8D") // Muted gray
	brandBg      = lipgloss.Color("#0D1117") // Deep dark background (virtual)

	titleStyle = lipgloss.NewStyle().
			Foreground(brandCyan).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(brandPurple)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(brandPurple).
			Padding(1, 2)

	activeItemStyle = lipgloss.NewStyle().
			Foreground(brandCyan).
			Bold(true)

	inactiveItemStyle = lipgloss.NewStyle().
				Foreground(brandGray)

	successStyle = lipgloss.NewStyle().
			Foreground(brandGreen)

	errorStyle = lipgloss.NewStyle().
			Foreground(brandRed)

	helpStyle = lipgloss.NewStyle().
			Foreground(brandGray).
			Italic(true)
)

const asciiBanner = `
 ███████╗ ██████╗ ███╗   ██╗███████╗
 ╚══███╔╝██╔═══██╗████╗  ██║██╔════╝
   ███╔╝ ██║   ██║██╔██╗ ██║█████╗  
  ███╔╝  ██║   ██║██║╚██╗██║██╔══╝  
 ███████╗╚██████╔╝██║ ╚████║███████╗
 ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝
          R E S T O R E
`

// Custom TUI Message definitions
type (
	logLineMsg      string
	execFinishedMsg struct {
		err error
	}
)

type tuiModel struct {
	step             int // 0 to 8
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

	applyFonts     bool // true = Yes, false = No
	fontCursor     int  // index in availableFonts
	availableFonts []string

	applyBg         bool // true = Yes, false = No
	bgChoiceCursor  int  // 0 = Predefined, 1 = Repo List, 2 = Custom Path
	wallpaperCursor int  // index in repoWallpapers
	repoWallpapers  []string

	// Defaults checkboxes state (Step 6)
	enableDocker      bool // Option A
	enableZshDefault  bool // Option B
	pinDiscord        bool // Option C
	configureKeyboard bool // Option D
	addArabic         bool // Option E
	configurePower    bool // Option F

	exportPrompt     bool // true = Yes, false = No
	importedSettings bool // true if settings were imported at Step 0
	exportDone       bool // true if final export step finished

	// Form text inputs
	gitNameInput  textinput.Model
	gitEmailInput textinput.Model
	customBgInput textinput.Model
	focusedInput  int // Cursor focus inside multi-option views (e.g. step 6)

	// Execution logs state
	applying bool
	finished bool
	err      error
	logChan  chan tea.Msg
	logLines []string

	// Responsive terminal size
	width   int
	height  int
	spinner spinner.Model
}
