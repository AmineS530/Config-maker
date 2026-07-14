package theme

import (
	"io"
	"os/exec"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.ApplyTheme {
		return nil
	}

	logger.Info("Applying theme settings: %s", cfg.ThemeName)
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", cfg.ThemeName).Run()
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.wm.preferences", "theme", cfg.ThemeName).Run()

	if cfg.ThemeMode == "1" {
		_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-dark").Run()
	} else {
		_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "default").Run()
	}

	_ = exec.Command("gsettings", "set", "org.gnome.shell.extensions.dash-to-dock", "extend-height", "false").Run()
	_ = exec.Command("gsettings", "set", "org.gnome.shell.extensions.dash-to-dock", "dock-fixed", "false").Run()
	logger.Success("System themes and dock panels set up.")
	return nil
}
