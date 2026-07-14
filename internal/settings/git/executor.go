package git

import (
	"io"
	"os/exec"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.ConfigureGit {
		return nil
	}

	logger.Info("Configuring Git settings...")
	_ = exec.Command("git", "config", "--global", "credential.helper", "store").Run()
	if cfg.GitName != "" {
		_ = exec.Command("git", "config", "--global", "user.name", cfg.GitName).Run()
	}
	if cfg.GitEmail != "" {
		_ = exec.Command("git", "config", "--global", "user.email", cfg.GitEmail).Run()
	}
	logger.Success("Git credentials configured globally.")
	return nil
}
