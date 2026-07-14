package keyboard

import (
	"io"
	"os/exec"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.ConfigureKeyboard {
		return nil
	}

	logger.Info("Configuring Gnome keyboard layouts...")
	var sources string
	if cfg.AddArabic {
		sources = "[('xkb', 'us'), ('xkb', 'fr'), ('xkb', 'ar')]"
	} else {
		sources = "[('xkb', 'us'), ('xkb', 'fr')]"
	}

	layoutCmd := exec.Command("gsettings", "set", "org.gnome.desktop.input-sources", "sources", sources)
	if err := layoutCmd.Run(); err != nil {
		logger.Warning("Failed to configure keyboard layouts: %v", err)
	} else {
		if cfg.AddArabic {
			logger.Success("Gnome keyboard layouts configured: US + FR + AR (Arabic).")
		} else {
			logger.Success("Gnome keyboard layouts configured: US + FR.")
		}
	}
	return nil
}
