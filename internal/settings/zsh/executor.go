package zsh

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"zonerestore/internal/utils"
)

//go:embed assets/.zshrc
var zshrcTemplate string

//go:embed assets/.p10k.zsh
var p10kTemplate string

type Alias struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Enabled bool   `json:"enabled"`
}

type TemplateData struct {
	CustomUsername string  `json:"custom_username"`
	Aliases        []Alias `json:"aliases"`
}

func Apply(cfg Config, data TemplateData, logger *utils.Logger, out io.Writer) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	if cfg.InstallOhMyZsh {
		zshDir := filepath.Join(homeDir, ".oh-my-zsh")
		if _, err := os.Stat(zshDir); os.IsNotExist(err) {
			logger.Info("Installing Oh-My-Zsh unattended...")
			installCmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | sh -s -- --unattended")
			if err := utils.RunCommandStream(installCmd, out); err != nil {
				logger.Warning("Oh-My-Zsh installation completed with warning: %v", err)
			} else {
				logger.Success("Oh-My-Zsh installed successfully.")
			}
		} else {
			logger.Info("Oh-My-Zsh is already installed. Skipping.")
		}
	}

	// Move Zsh configs
	logger.Info("Generating premade Zsh and Powerlevel10k configurations from embedded templates...")
	p10kDst := filepath.Join(homeDir, ".p10k.zsh")
	if err := writeTemplate(p10kTemplate, p10kDst, data); err != nil {
		logger.Warning("Failed to generate .p10k.zsh: %v", err)
	} else {
		logger.Success("Generated .p10k.zsh successfully.")
	}

	zshrcDst := filepath.Join(homeDir, ".zshrc")
	if err := writeTemplate(zshrcTemplate, zshrcDst, data); err != nil {
		logger.Warning("Failed to generate .zshrc: %v", err)
	} else {
		logger.Success("Generated .zshrc successfully.")
	}

	// Set up Powerlevel10k theme
	p10kThemeDir := filepath.Join(homeDir, "powerlevel10k")
	if _, err := os.Stat(p10kThemeDir); os.IsNotExist(err) {
		logger.Info("Cloning Powerlevel10k repository...")
		p10kCloneCmd := exec.Command("git", "clone", "--quiet", "--depth=1", "https://github.com/romkatv/powerlevel10k.git", p10kThemeDir)
		if err := utils.RunCommandStream(p10kCloneCmd, out); err == nil {
			logger.Success("Powerlevel10k cloned successfully.")
			_ = utils.AppendToFile(zshrcDst, "source ~/powerlevel10k/powerlevel10k.zsh-theme\n")
		} else {
			logger.Warning("Failed to clone Powerlevel10k theme repository: %v", err)
		}
	} else {
		logger.Info("Powerlevel10k directory already exists. Skipping.")
	}

	return nil
}

func writeTemplate(tmplContent, destPath string, data TemplateData) error {
	tmpl, err := template.New("config").Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(destPath, buf.Bytes(), 0o644)
}
