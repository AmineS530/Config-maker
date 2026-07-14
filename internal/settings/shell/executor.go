package shell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	bashrcPath := filepath.Join(homeDir, ".bashrc")
	if cfg.EnableZshDefault {
		logger.Info("Configuring Zsh as default shell in .bashrc...")
		if err := utils.AppendZshToBashrc(bashrcPath); err != nil {
			logger.Error("Failed to switch default shell to Zsh: %v", err)
		} else {
			logger.Success("Zsh default shell configuration added to .bashrc.")
		}
	} else {
		logger.Info("Ensuring Zsh is not default shell in .bashrc...")
		if err := utils.RemoveZshFromBashrc(bashrcPath); err != nil {
			logger.Error("Failed to restore default shell: %v", err)
		} else {
			logger.Success("Zsh default shell configuration removed from .bashrc.")
		}
	}
	return nil
}
