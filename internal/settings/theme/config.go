package theme

type Config struct {
	ApplyTheme bool   `json:"apply_theme"`
	ThemeMode  string `json:"theme_mode"` // "1" = dark, "2" = light
	ThemeName  string `json:"theme_name"`
}
