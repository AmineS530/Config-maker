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

	// Copy all fonts from themes/fonts/ (both flat files and subdirectories)
	fontDir := filepath.Join(themesRoot, "fonts")
	var fontsCopied int
	if info, err := os.Stat(fontDir); err == nil && info.IsDir() {
		// Walk the whole fonts directory — handles both flat and subdir layouts
		_ = filepath.Walk(fontDir, func(path string, info os.FileInfo, err error) error {
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
		logger.Success("Installed %d font files to ~/.local/share/fonts.", fontsCopied)
	} else {
		logger.Warning("No font files found in %s to copy.", fontDir)
	}

	logger.Info("Refreshing system font cache...")
	_ = exec.Command("fc-cache", "-f", fontsTargetDir).Run()

	// ── Terminal font: always MesloLGS NF — hardcoded, not user-overridable ──
	terminalFont := fmt.Sprintf("%s %d", TerminalFont, TerminalFontSize)
	logger.Info("Setting terminal font: %s (locked)", terminalFont)

	// Apply terminal font to all GNOME Terminal profiles
	if profileIDs, err := utils.GetGnomeTerminalProfiles(); err == nil && len(profileIDs) > 0 {
		logger.Info("Configuring terminal profile fonts...")
		for _, pID := range profileIDs {
			_ = exec.Command("dconf", "write",
				fmt.Sprintf("/org/gnome/terminal/legacy/profiles:/:%s/font", pID),
				fmt.Sprintf("'%s'", terminalFont),
			).Run()
			_ = exec.Command("dconf", "write",
				fmt.Sprintf("/org/gnome/terminal/legacy/profiles:/:%s/use-system-font", pID),
				"false",
			).Run()
		}
		logger.Success("Terminal font set to %s in all GNOME Terminal profiles.", terminalFont)
	}

	// Also set system monospace font to MesloLGS so editors/IDEs pick it up
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", terminalFont).Run()

	// ── Display font: user-selected, applied to GNOME interface ──
	displayFont := cfg.DisplayFontName
	if displayFont == "" {
		// Sensible default: Ubuntu 11 (GNOME default on Ubuntu)
		displayFont = "Ubuntu 11"
	} else {
		displayFont = displayFont + " 11"
	}
	logger.Info("Setting GNOME display font: %s", displayFont)
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "font-name", displayFont).Run()
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "document-font-name", displayFont).Run()

	logger.Success("Font configuration applied.")
	return nil
}
