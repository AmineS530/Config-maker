package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"zonerestore/internal/settings/dock"
	"zonerestore/internal/settings/docker"
	"zonerestore/internal/settings/fonts"
	"zonerestore/internal/settings/git"
	"zonerestore/internal/settings/keyboard"
	"zonerestore/internal/settings/power"
	"zonerestore/internal/settings/shell"
	"zonerestore/internal/settings/theme"
	"zonerestore/internal/settings/wallpaper"
	"zonerestore/internal/settings/zsh"
)

// UserConfig holds the user choices for the configuration wizard composed of context-specific modules.
type UserConfig struct {
	Zsh            zsh.Config       `json:"zsh"`
	Git            git.Config       `json:"git"`
	Theme          theme.Config     `json:"theme"`
	Fonts          fonts.Config     `json:"fonts"`
	Wallpaper      wallpaper.Config `json:"wallpaper"`
	Docker         docker.Config    `json:"docker"`
	Dock           dock.Config      `json:"dock"`
	Keyboard       keyboard.Config  `json:"keyboard"`
	Power          power.Config     `json:"power"`
	Shell          shell.Config     `json:"shell"`
	CustomUsername string           `json:"custom_username"`
	Aliases        []zsh.Alias      `json:"aliases"`
}

// DefaultConfig returns a pre-populated default configuration.
func DefaultConfig() UserConfig {
	return UserConfig{
		Zsh: zsh.Config{
			InstallOhMyZsh: true,
		},
		Git: git.Config{
			ConfigureGit: true,
			GitName:      "",
			GitEmail:     "",
		},
		Theme: theme.Config{
			ApplyTheme: true,
			ThemeMode:  "1",
			ThemeName:  "",
		},
		Fonts: fonts.Config{
			ConfigureFonts: true,
			FontName:       "MesloLGS NF",
		},
		Wallpaper: wallpaper.Config{
			ApplyBackground: true,
			BackgroundImage: "",
		},
		Docker: docker.Config{
			EnableDocker: true,
		},
		Dock: dock.Config{
			PinDiscord: true,
		},
		Keyboard: keyboard.Config{
			ConfigureKeyboard: true,
			AddArabic:         false,
		},
		Power: power.Config{
			ConfigurePower: true,
		},
		Shell: shell.Config{
			EnableZshDefault: false,
		},
		CustomUsername: GetDefaultUsername(),
		Aliases: []zsh.Alias{
			{
				Name:    "quickpush",
				Command: `gofmt -w . && gaa && gc -m "quick_add_commit_push_alias" && gp`,
				Enabled: false,
			},
		},
	}
}

// GetDefaultUsername returns the current logged-in username or "user" as fallback.
func GetDefaultUsername() string {
	u, err := user.Current()
	if err == nil && u.Username != "" {
		return u.Username
	}
	if uEnv := os.Getenv("USER"); uEnv != "" {
		return uEnv
	}
	return "user"
}

// LoadConfig attempts to read a saved configuration from ~/.config/zonerestore/config.json.
// If it fails or the file doesn't exist, it falls back to DefaultConfig().
func LoadConfig() UserConfig {
	cfg := DefaultConfig()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	configFilePath := filepath.Join(homeDir, ".config", "zonerestore", "config.json")
	file, err := os.Open(configFilePath)
	if err != nil {
		return cfg // doesn't exist or can't open
	}
	defer file.Close()

	var loaded UserConfig
	err = json.NewDecoder(file).Decode(&loaded)
	if err == nil {
		// Fill in empty subfields if missing from older versions of config
		return loaded
	}
	return cfg
}

// SaveConfig writes the given UserConfig to ~/.config/zonerestore/config.json.
func SaveConfig(cfg UserConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(homeDir, ".config", "zonerestore")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	configFilePath := filepath.Join(configDir, "config.json")
	configData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, configData, 0644)
}
