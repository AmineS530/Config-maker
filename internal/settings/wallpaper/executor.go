package wallpaper

import (
	"io"
	"os/exec"
	"path/filepath"

	"zonerestore/internal/themes"
	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.ApplyBackground {
		return nil
	}

	bgPath := cfg.BackgroundImage
	if bgPath == "" {
		bgPath = themes.WallpaperPath(themes.Root(), "Background.jpeg")
	} else if !filepath.IsAbs(bgPath) {
		bgPath = themes.WallpaperPath(themes.Root(), bgPath)
	}

	logger.Info("Applying desktop background: %s", bgPath)
	fileURL := "file://" + bgPath
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", fileURL).Run()
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", fileURL).Run()
	logger.Success("Desktop wallpaper applied.")
	return nil
}
