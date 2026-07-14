package dock

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"zonerestore/internal/utils"
)

func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
	if !cfg.PinDiscord {
		return nil
	}

	logger.Info("Configuring Gnome Shell favorites in dock...")

	// 1. Read current favorites
	cmd := exec.Command("gsettings", "get", "org.gnome.shell", "favorite-apps")
	outputBytes, err := cmd.Output()
	var current []string
	if err == nil {
		currentStr := strings.TrimSpace(string(outputBytes))
		// Remove brackets
		currentStr = strings.TrimPrefix(currentStr, "[")
		currentStr = strings.TrimSuffix(currentStr, "]")

		if currentStr != "" {
			parts := strings.Split(currentStr, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				p = strings.Trim(p, "'\"")
				if p != "" {
					current = append(current, p)
				}
			}
		}
	}

	// 2. Check if discord is already pinned
	alreadyPinned := false
	for _, app := range current {
		if app == "discord.desktop" {
			alreadyPinned = true
			break
		}
	}

	if alreadyPinned {
		logger.Success("Discord is already pinned to favorites.")
		return nil
	}

	// 3. Append discord
	current = append(current, "discord.desktop")

	// 4. Construct back to gsettings format
	var sb strings.Builder
	sb.WriteString("[")
	for i, app := range current {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("'%s'", app))
	}
	sb.WriteString("]")

	// 5. Apply new list
	setCmd := exec.Command("gsettings", "set", "org.gnome.shell", "favorite-apps", sb.String())
	if err := setCmd.Run(); err != nil {
		logger.Warning("Failed to pin Discord to favorites: %v", err)
		return nil
	}

	logger.Success("Discord pinned to favorites successfully.")
	return nil
}
