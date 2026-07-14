package fonts

type Config struct {
	ConfigureFonts  bool   `json:"configure_fonts"`
	FontName        string `json:"font_name"`         // Terminal / monospace font — always MesloLGS NF
	DisplayFontName string `json:"display_font_name"` // GNOME interface / display font (optional, user-selectable)
}

// TerminalFont is the hardcoded monospace font always applied to GNOME Terminal.
// This is intentionally not user-configurable to maintain a consistent developer experience.
const TerminalFont = "MesloLGS NF Regular"
const TerminalFontSize = 12
