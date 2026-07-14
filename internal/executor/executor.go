package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"zonerestore/internal/config"
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
	"zonerestore/internal/utils"
)

// ApplyConfig runs the orchestration of the setup steps natively in Go.
func ApplyConfig(cfg config.UserConfig, exportSettings bool, out io.Writer) error {
	logger := &utils.Logger{Out: out}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	logger.Info("Starting configuration application natively in Go...")

	// 1. Oh-My-Zsh & Zsh config setup
	// Convert config aliases to settings aliases
	var aliases []zsh.Alias
	for _, a := range cfg.Aliases {
		aliases = append(aliases, zsh.Alias(a))
	}
	zshData := zsh.TemplateData{
		CustomUsername: cfg.CustomUsername,
		Aliases:        aliases,
	}
	if err := zsh.Apply(cfg.Zsh, zshData, logger, out); err != nil {
		logger.Warning("Oh-My-Zsh setup completed with warning: %v", err)
	}

	// 2. Clone the Config-maker repository
	destDir := filepath.Join(homeDir, "ZoneRestore")
	logger.Info("Setting up repository destination: %s", destDir)

	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		logger.Warning("Directory already exists. Removing older installation...")
		if err := os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("failed to clear destination directory: %w", err)
		}
	}

	logger.Info("Cloning configuration repository (shallow)...")
	cloneCmd := exec.Command("git", "clone", "--quiet", "--depth=1", "https://github.com/AmineS530/Config-maker.git", destDir)
	if err := utils.RunCommandStream(cloneCmd, out); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	logger.Success("Repository cloned successfully.")

	// 3. System Gsettings (layout & power settings)
	if err := keyboard.Apply(cfg.Keyboard, logger, out); err != nil {
		logger.Warning("Keyboard configuration completed with warning: %v", err)
	}
	if err := power.Apply(cfg.Power, logger, out); err != nil {
		logger.Warning("Power configuration completed with warning: %v", err)
	}

	// 4. Git Global Credentials
	if err := git.Apply(cfg.Git, logger, out); err != nil {
		logger.Warning("Git configuration completed with warning: %v", err)
	}

	// 5. System Gsettings Theme
	if err := theme.Apply(cfg.Theme, logger, out); err != nil {
		logger.Warning("Theme application completed with warning: %v", err)
	}

	// 6. Wallpaper Application
	if err := wallpaper.Apply(cfg.Wallpaper, logger, out); err != nil {
		logger.Warning("Wallpaper application completed with warning: %v", err)
	}

	// 7. Custom Fonts Installation
	if err := fonts.Apply(cfg.Fonts, logger, out); err != nil {
		logger.Warning("Fonts installation completed with warning: %v", err)
	}

	// 8. Docker Rootless
	if err := docker.Apply(cfg.Docker, logger, out); err != nil {
		logger.Warning("Docker setup completed with warning: %v", err)
	}

	// 9. Dock Favorites (Pin Discord)
	if err := dock.Apply(cfg.Dock, logger, out); err != nil {
		logger.Warning("Dock favorites setup completed with warning: %v", err)
	}

	// 10. Default Shell Switch
	if err := shell.Apply(cfg.Shell, logger, out); err != nil {
		logger.Warning("Default shell configuration completed with warning: %v", err)
	}

	if exportSettings {
		logger.Info("Exporting configuration settings...")
		if err := config.SaveConfig(cfg); err != nil {
			logger.Error("Failed to save settings: %v", err)
		} else {
			logger.Success("Settings exported successfully.")
		}
	}

	logger.Success("All configurations applied successfully!")
	return nil
}

// FinishSetup performs the terminal closing and reloading behavior.
func FinishSetup(out io.Writer) {
	logger := &utils.Logger{Out: out}
	logger.Info("Finishing setup. Close other terminals and prompt reload...")

	// 1. Pop Zenity info dialog
	zenityText := `<span size="x-large">Finished.</span>\n\nPlease restart your terminal to apply.\n` +
		`In case of encountering a problem, send a PM to a.sadik on Discord.\n\n` +
		`Terminal will attempt to re-open after clicking\n\n<b>Happy coding :)</b>`

	zenityCmd := exec.Command(
		"zenity",
		"--info",
		"--text="+zenityText,
		"--title=Setup is done",
		"--ok-label=yessir",
	)
	_ = zenityCmd.Run()

	// 2. Restart terminal (run gnome-terminal in background before killing this instance)
	reopenCmd := exec.Command("gnome-terminal")
	_ = reopenCmd.Start()

	// 3. Kill GNOME Terminal instances
	logger.Warning("Closing all running GNOME Terminal instances...")
	killCmd := exec.Command("killall", "gnome-terminal-server")
	_ = killCmd.Run()
}
