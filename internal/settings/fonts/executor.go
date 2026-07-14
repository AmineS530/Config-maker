package fonts

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"zonerestore/internal/themes"
	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.ConfigureFonts {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	logger.Info("Installing custom fonts...")
	fontsTargetDir := filepath.Join(homeDir, ".local/share/fonts")
	_ = os.MkdirAll(fontsTargetDir, 0o755)

	themesRoot := themes.Root()
	var fontSourceDir string
	if cfg.FontName != "" {
		fontSourceDir = filepath.Join(themesRoot, "fonts", cfg.FontName)
	} else {
		fontSourceDir = filepath.Join(themesRoot, "fonts")
	}

	var fontsCopied int
	if info, err := os.Stat(fontSourceDir); err == nil && info.IsDir() {
		_ = filepath.Walk(fontSourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".ttf" || ext == ".otf" || ext == ".woff" || ext == ".woff2" {
				dstPath := filepath.Join(fontsTargetDir, info.Name())
				if err := utils.CopyFile(path, dstPath); err == nil {
					fontsCopied++
				}
			}
			return nil
		})
	}

	if fontsCopied > 0 {
		logger.Success("Installed %d custom fonts to local storage.", fontsCopied)
	} else {
		logger.Warning("No custom fonts found in %s to copy.", fontSourceDir)
	}

	logger.Info("Refreshing system font cache...")
	_ = exec.Command("fc-cache", "-f", "-v", fontsTargetDir).Run()

	// Apply Fonts
	logger.Info("Applying Gnome interface fonts...")
	monoFont := "MesloLGS NF Regular 12"
	if cfg.FontName != "" {
		monoFont = cfg.FontName + " 12"
	}
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "font-name", "Ubuntu 11").Run()
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", monoFont).Run()

	// Set Gnome terminal default profile font
	if profileIDs, err := utils.GetGnomeTerminalProfiles(); err == nil && len(profileIDs) > 0 {
		logger.Info("Configuring terminal profile fonts...")
		for _, pID := range profileIDs {
			_ = exec.Command("dconf", "write", fmt.Sprintf("/org/gnome/terminal/legacy/profiles:/:%s/font", pID), fmt.Sprintf("'%s'", monoFont)).Run()
			_ = exec.Command("dconf", "write", fmt.Sprintf("/org/gnome/terminal/legacy/profiles:/:%s/use-system-font", pID), "false").Run()
		}
		logger.Success("Monospace font configured in gnome-terminal profiles.")
	}

	return nil
}
