package power

import (
	"io"
	"os/exec"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	// if !cfg.ConfigurePower {
	// 	return nil
	// }

	logger.Info("Configuring Gnome power management settings...")
	powerCmd := exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-timeout", "1800")
	if err := powerCmd.Run(); err != nil {
		logger.Warning("Failed to configure power settings: %v", err)
	} else {
		logger.Success("Gnome sleep timeout set to 0.5 hours.")
	}
	return nil
}
