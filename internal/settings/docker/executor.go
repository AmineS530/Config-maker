package docker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.EnableDocker {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	logger.Info("Setting up Docker in rootless mode...")
	dockerCmd := exec.Command("sh", "-c", "curl -fsSL https://get.docker.com/rootless | sh")
	dockerCmd.Env = append(os.Environ(), "PATH="+filepath.Join(homeDir, "bin")+":"+os.Getenv("PATH"))
	if err := utils.RunCommandStream(dockerCmd, out); err != nil {
		logger.Warning("Docker rootless installation completed with warning: %v", err)
	} else {
		logger.Success("Docker rootless installed successfully.")
		fmt.Fprintf(out, "\n%s[DOCKER CONFIG]%s To run docker, copy/paste these variables or restart terminal:\n", utils.Green, utils.Reset)
		fmt.Fprintf(out, "export PATH=%s/bin:$PATH\n", homeDir)
		fmt.Fprintf(out, "export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock\n\n")
	}
	return nil
}
