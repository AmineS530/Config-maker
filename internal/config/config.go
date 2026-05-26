package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UserConfig holds the user choices for the configuration wizard.
type UserConfig struct {
	InstallOhMyZsh   bool   `json:"install_oh_my_zsh"`
	ConfigureGit     bool   `json:"configure_git"`
	GitName          string `json:"git_name"`
	GitEmail         string `json:"git_email"`
	ApplyTheme       bool   `json:"apply_theme"`
	ThemeMode        string `json:"theme_mode"` // "1" = dark, "2" = light
	ThemeName        string `json:"theme_name"`
	ApplyBackground  bool   `json:"apply_background"`
	BackgroundImage  string `json:"background_image"`
	EnableDocker     bool   `json:"enable_docker"`
	EnableZshDefault bool   `json:"enable_zsh_default"`
}

// DefaultConfig returns a pre-populated default configuration.
func DefaultConfig() UserConfig {
	return UserConfig{
		InstallOhMyZsh:   true,
		ConfigureGit:     true,
		ApplyTheme:       true,
		ThemeMode:        "1",
		ApplyBackground:  true,
		EnableDocker:     false,
		EnableZshDefault: true,
	}
}

// LoadConfig attempts to read a saved configuration from ~/.config/config-maker/config.json.
// If it fails or the file doesn't exist, it falls back to DefaultConfig().
func LoadConfig() UserConfig {
	cfg := DefaultConfig()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	configFilePath := filepath.Join(homeDir, ".config", "config-maker", "config.json")
	file, err := os.Open(configFilePath)
	if err != nil {
		return cfg // doesn't exist or can't open
	}
	defer file.Close()

	var loaded UserConfig
	err = json.NewDecoder(file).Decode(&loaded)
	if err == nil {
		return loaded
	}
	return cfg
}

// SaveConfig writes the given UserConfig to ~/.config/config-maker/config.json.
func SaveConfig(cfg UserConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(homeDir, ".config", "config-maker")
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
