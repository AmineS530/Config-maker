package executor

import (
	"bufio"
	"bytes"
	"config-maker/internal/config"
	"config-maker/internal/utils"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed assets/.zshrc
var zshrcTemplate string

//go:embed assets/.p10k.zsh
var p10kTemplate string

// ApplyConfig runs the orchestration of the setup steps natively in Go.
func ApplyConfig(cfg config.UserConfig, exportSettings bool, out io.Writer) error {
	logger := &utils.Logger{Out: out}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	logger.Info("Starting configuration application natively in Go...")

	// 1. Install Oh-My-Zsh if requested
	if cfg.InstallOhMyZsh {
		zshDir := filepath.Join(homeDir, ".oh-my-zsh")
		if _, err := os.Stat(zshDir); os.IsNotExist(err) {
			logger.Info("Installing Oh-My-Zsh unattended...")
			installCmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | sh -s -- --unattended")
			if err := runCommandStream(installCmd, out); err != nil {
				logger.Warning("Oh-My-Zsh installation completed with warning: %v", err)
			} else {
				logger.Success("Oh-My-Zsh installed successfully.")
			}
		} else {
			logger.Info("Oh-My-Zsh is already installed. Skipping.")
		}
	}

	// 2. Clone the Config-maker repository
	destDir := filepath.Join(homeDir, "Zone01_Desk_cfg")
	logger.Info("Setting up repository destination: %s", destDir)

	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		logger.Warning("Directory already exists. Removing older installation...")
		if err := os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("failed to clear destination directory: %w", err)
		}
	}

	logger.Info("Cloning configuration repository (shallow)...")
	cloneCmd := exec.Command("git", "clone", "--quiet", "--depth=1", "https://github.com/AmineS530/Config-maker.git", destDir)
	if err := runCommandStream(cloneCmd, out); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	logger.Success("Repository cloned successfully.")

	// 3. Move premade Zsh configurations
	logger.Info("Generating premade Zsh and Powerlevel10k configurations from embedded templates...")
	p10kDst := filepath.Join(homeDir, ".p10k.zsh")
	if err := writeTemplate(p10kTemplate, p10kDst, cfg); err != nil {
		logger.Warning("Failed to generate .p10k.zsh: %v", err)
	} else {
		logger.Success("Generated .p10k.zsh successfully.")
	}

	zshrcDst := filepath.Join(homeDir, ".zshrc")
	if err := writeTemplate(zshrcTemplate, zshrcDst, cfg); err != nil {
		logger.Warning("Failed to generate .zshrc: %v", err)
	} else {
		logger.Success("Generated .zshrc successfully.")
	}

	// 4. Set up Powerlevel10k theme
	p10kThemeDir := filepath.Join(homeDir, "powerlevel10k")
	if _, err := os.Stat(p10kThemeDir); os.IsNotExist(err) {
		logger.Info("Cloning Powerlevel10k repository...")
		p10kCloneCmd := exec.Command("git", "clone", "--quiet", "--depth=1", "https://github.com/romkatv/powerlevel10k.git", p10kThemeDir)
		if err := runCommandStream(p10kCloneCmd, out); err == nil {
			logger.Success("Powerlevel10k cloned successfully.")
			_ = appendToFile(zshrcDst, "source ~/powerlevel10k/powerlevel10k.zsh-theme\n")
		} else {
			logger.Warning("Failed to clone Powerlevel10k theme repository: %v", err)
		}
	} else {
		logger.Info("Powerlevel10k directory already exists. Skipping.")
	}

	// 5. System Gsettings (layout & power settings)
	if cfg.ConfigureKeyboard {
		logger.Info("Configuring Gnome keyboard layouts...")
		layoutCmd := exec.Command("gsettings", "set", "org.gnome.desktop.input-sources", "sources", "[('xkb', 'us'), ('xkb', 'fr')]")
		if err := layoutCmd.Run(); err != nil {
			logger.Warning("Failed to configure keyboard layouts: %v", err)
		} else {
			logger.Success("Gnome keyboard layouts configured.")
		}
	}

	if cfg.ConfigurePower {
		logger.Info("Configuring Gnome power settings...")
		powerCmd := exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-timeout", "5400")
		if err := powerCmd.Run(); err != nil {
			logger.Warning("Failed to configure power settings: %v", err)
		} else {
			logger.Success("Gnome sleep timeout set to 1.5 hours.")
		}
	}

	// 6. Configure Git credentials
	if cfg.ConfigureGit {
		logger.Info("Configuring Git settings...")
		_ = exec.Command("git", "config", "--global", "credential.helper", "store").Run()
		if cfg.GitName != "" {
			_ = exec.Command("git", "config", "--global", "user.name", cfg.GitName).Run()
		}
		if cfg.GitEmail != "" {
			_ = exec.Command("git", "config", "--global", "user.email", cfg.GitEmail).Run()
		}
		logger.Success("Git credentials configured globally.")
	}

	// 7. Apply session theme
	if cfg.ApplyTheme {
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
	}

	// 8. Apply Background Image
	if cfg.ApplyBackground {
		bgPath := cfg.BackgroundImage
		if bgPath == "" {
			bgPath = filepath.Join(destDir, "Background.jpeg")
		}
		logger.Info("Applying desktop background: %s", bgPath)
		fileURL := "file://" + bgPath
		_ = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", fileURL).Run()
		_ = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", fileURL).Run()
		logger.Success("Desktop wallpaper applied.")
	}

	// 9. Font Installation & Terminal Font Setup
	if cfg.ConfigureFonts {
		logger.Info("Installing custom fonts...")
		fontsTargetDir := filepath.Join(homeDir, ".local", "share", "fonts")
		_ = os.MkdirAll(fontsTargetDir, 0755)

		displayFontFile := "MPLUS1p-Regular.ttf"
		terminalFontFile := "MesloLGS NF Regular.ttf"

		err1 := copyFile(filepath.Join(destDir, "fonts", displayFontFile), filepath.Join(fontsTargetDir, displayFontFile))
		err2 := copyFile(filepath.Join(destDir, "fonts", terminalFontFile), filepath.Join(fontsTargetDir, terminalFontFile))

		if err1 != nil || err2 != nil {
			logger.Warning("Font copy completed with warnings. (Display: %v, Monospace: %v)", err1, err2)
		} else {
			logger.Success("Custom fonts copied to local storage.")
		}

		logger.Info("Refreshing system font cache...")
		_ = exec.Command("fc-cache", "-f", "-v", fontsTargetDir).Run()

		// Apply Fonts
		logger.Info("Applying Gnome interface fonts...")
		_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "font-name", "MPLUS1p-Regular 12").Run()
		_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", "MesloLGS NF Regular 12").Run()

		// Set Gnome terminal default profile font
		if profileIDs, err := getGnomeTerminalProfiles(); err == nil && len(profileIDs) > 0 {
			logger.Info("Configuring terminal profile fonts...")
			for _, pID := range profileIDs {
				_ = exec.Command("dconf", "write", fmt.Sprintf("/org/gnome/terminal/legacy/profiles:/:%s/font", pID), "'MesloLGS NF Regular 12'").Run()
				_ = exec.Command("dconf", "write", fmt.Sprintf("/org/gnome/terminal/legacy/profiles:/:%s/use-system-font", pID), "false").Run()
			}
			logger.Success("Monospace font configured in gnome-terminal profiles.")
		}
	}

	// 10. Set up Docker Rootless if requested
	if cfg.EnableDocker {
		logger.Info("Setting up Docker in rootless mode...")
		dockerCmd := exec.Command("sh", "-c", "curl -fsSL https://get.docker.com/rootless | sh")
		dockerCmd.Env = append(os.Environ(), "PATH="+filepath.Join(homeDir, "bin")+":"+os.Getenv("PATH"))
		if err := runCommandStream(dockerCmd, out); err != nil {
			logger.Warning("Docker rootless installation completed with warning: %v", err)
		} else {
			logger.Success("Docker rootless installed successfully.")
			fmt.Fprintf(out, "\n%s[DOCKER CONFIG]%s To run docker, copy/paste these variables or restart terminal:\n", utils.Green, utils.Reset)
			fmt.Fprintf(out, "export PATH=%s/bin:$PATH\n", homeDir)
			fmt.Fprintf(out, "export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock\n\n")
		}
	}

	// 11. Handle Zsh Shell Switching (enabling/disabling)
	bashrcPath := filepath.Join(homeDir, ".bashrc")
	if cfg.EnableZshDefault {
		logger.Info("Configuring Zsh as default shell in .bashrc...")
		if err := appendZshToBashrc(bashrcPath); err != nil {
			logger.Error("Failed to switch default shell to Zsh: %v", err)
		} else {
			logger.Success("Zsh default shell configuration added to .bashrc.")
		}
	} else {
		logger.Info("Ensuring Zsh is not default shell in .bashrc...")
		if err := removeZshFromBashrc(bashrcPath); err != nil {
			logger.Error("Failed to remove Zsh default shell configuration: %v", err)
		} else {
			logger.Success("Zsh default shell configuration removed from .bashrc.")
		}
	}

	// 12. Export selections to $HOME/.config/config-maker/config.json if requested
	if exportSettings {
		if err := config.SaveConfig(cfg); err != nil {
			logger.Error("Failed to save config: %v", err)
		} else {
			configFilePath := filepath.Join(homeDir, ".config", "config-maker", "config.json")
			logger.Success("Exported selections to config file: %s", configFilePath)
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

	zenityCmd := exec.Command("zenity",
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

// runCommandStream executes a command and prints its stdout/stderr live to the writer.
func runCommandStream(cmd *exec.Cmd, out io.Writer) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Read stdout and stderr concurrently
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Fprintln(out, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Fprintln(out, scanner.Text())
		}
	}()

	return cmd.Wait()
}

// copyFileIfExists performs file copies if source exists.
func copyFileIfExists(src, dst string, logger *utils.Logger) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil // skip silently if doesn't exist
	}
	return copyFile(src, dst)
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// appendToFile appends a string to the specified file.
func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// getGnomeTerminalProfiles returns active gnome terminal legacy profiles UUID strings.
func getGnomeTerminalProfiles() ([]string, error) {
	outBytes, err := exec.Command("dconf", "list", "/org/gnome/terminal/legacy/profiles:/").Output()
	if err != nil {
		return nil, err
	}

	var profiles []string
	lines := strings.Split(string(outBytes), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		// Strip trailing slashes
		trimmed = strings.TrimSuffix(trimmed, "/")
		profiles = append(profiles, trimmed)
	}
	return profiles, nil
}

// appendZshToBashrc appends default Zsh shell settings to .bashrc.
func appendZshToBashrc(bashrcPath string) error {
	content, err := os.ReadFile(bashrcPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	zshLines := "SHELL=/bin/zsh\nexec /bin/zsh -l\n"
	if strings.Contains(string(content), "SHELL=/bin/zsh") {
		return nil // already exists
	}

	return appendToFile(bashrcPath, zshLines)
}

// removeZshFromBashrc removes Zsh shell settings from .bashrc.
func removeZshFromBashrc(bashrcPath string) error {
	content, err := os.ReadFile(bashrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "SHELL=/bin/zsh" || trimmed == "exec /bin/zsh -l" {
			continue
		}
		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	if string(content) == newContent {
		return nil
	}

	return os.WriteFile(bashrcPath, []byte(newContent), 0644)
}

// writeTemplate compiles a text template with cfg config and writes it to destPath.
func writeTemplate(tmplContent, destPath string, cfg config.UserConfig) error {
	tmpl, err := template.New("config").Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(destPath, buf.Bytes(), 0644)
}
