// SHARED UTILITY: This package is shared between the CLI interface (internal/cli)
// and the Web server (internal/web). Modifying public APIs will impact both contexts.

package themes

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	RepoURL      = "https://github.com/AmineS530/ZoneRestoreThemes"
	RelativePath = ".local/share/zonerestore/themes"
)

// Root returns the absolute path to the themes directory without cloning.
func Root() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, RelativePath)
}

// EnsureThemes clones or updates the ZoneRestoreThemes repo.
// Returns the root path of the themes directory.
func EnsureThemes(ctx context.Context, w io.Writer) (string, error) {
	dest := Root()
	if dest == "" {
		return "", fmt.Errorf("could not determine user home directory")
	}

	gitDir := filepath.Join(dest, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		// Clone it
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return "", fmt.Errorf("failed to create base directory: %w", err)
		}
		fmt.Fprintf(w, "Cloning themes repository from %s...\n", RepoURL)
		cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", RepoURL, dest)
		cmd.Stdout = w
		cmd.Stderr = w
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to clone themes repository: %w", err)
		}
		fmt.Fprintln(w, "Themes repository cloned successfully.")
	} else {
		// Update it
		fmt.Fprintln(w, "Updating themes repository...")
		cmd := exec.CommandContext(ctx, "git", "-C", dest, "pull", "--ff-only")
		cmd.Stdout = w
		cmd.Stderr = w
		if err := cmd.Run(); err != nil {
			// If pull fails (e.g. offline), log a warning but proceed with cached version
			fmt.Fprintf(w, "Warning: failed to update themes repository: %v. Using cached files.\n", err)
		} else {
			fmt.Fprintln(w, "Themes repository updated successfully.")
		}
	}

	return dest, nil
}

// ListWallpapers returns all image filenames in <root>/wallpapers/.
func ListWallpapers(root string) []string {
	if root == "" {
		return nil
	}
	wpDir := filepath.Join(root, "wallpapers")
	entries, err := os.ReadDir(wpDir)
	if err != nil {
		return nil
	}
	var wallpapers []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		nameLower := strings.ToLower(name)
		if strings.HasSuffix(nameLower, ".png") || strings.HasSuffix(nameLower, ".jpg") || strings.HasSuffix(nameLower, ".jpeg") {
			wallpapers = append(wallpapers, name)
		}
	}
	return wallpapers
}

// ListFonts returns all font family directories in <root>/fonts/.
func ListFonts(root string) []string {
	if root == "" {
		return nil
	}
	fontDir := filepath.Join(root, "fonts")
	entries, err := os.ReadDir(fontDir)
	if err != nil {
		return nil
	}
	var fonts []string
	for _, entry := range entries {
		if entry.IsDir() {
			fonts = append(fonts, entry.Name())
		}
	}
	return fonts
}

// WallpaperPath returns the absolute path to a given wallpaper file.
func WallpaperPath(root, name string) string {
	if root == "" || name == "" {
		return ""
	}
	return filepath.Join(root, "wallpapers", name)
}

// FontPath returns the absolute path to a given font directory.
func FontPath(root, name string) string {
	if root == "" || name == "" {
		return ""
	}
	return filepath.Join(root, "fonts", name)
}

// IsAvailable returns true if the themes directory exists and has content.
func IsAvailable(root string) bool {
	if root == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(root, ".git"))
	return err == nil
}
