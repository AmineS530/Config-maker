package web

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"zonerestore/internal/themes"
)

// isFontExt returns true if the extension is a recognized font format.
func isFontExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".ttf", ".otf", ".woff", ".woff2":
		return true
	}
	return false
}

// HandleResources returns available theme, wallpaper, and font options as JSON.
// Fonts are scanned from both subdirectories AND flat files in themes/fonts/.
func HandleResources(w http.ResponseWriter, r *http.Request) {
	darkThemes := getSystemThemes("1")
	lightThemes := getSystemThemes("2")

	themesRoot := themes.Root()
	wallpapers := themes.ListWallpapers(themesRoot)

	// ── Fonts: support both flat files AND subdirectory layout ──
	var fontDetails []FontDetail
	fontDir := filepath.Join(themesRoot, "fonts")
	entries, err := os.ReadDir(fontDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				// Subdirectory layout: fonts/<FamilyName>/*.ttf
				subDir := filepath.Join(fontDir, entry.Name())
				subEntries, err := os.ReadDir(subDir)
				if err != nil {
					continue
				}
				var files []string
				for _, sub := range subEntries {
					if !sub.IsDir() && isFontExt(filepath.Ext(sub.Name())) {
						files = append(files, sub.Name())
					}
				}
				if len(files) > 0 {
					fontDetails = append(fontDetails, FontDetail{
						Name:  entry.Name(),
						Files: files,
					})
				}
			} else {
				// Flat file layout: fonts/<FontFile.ttf>
				if isFontExt(filepath.Ext(entry.Name())) {
					// Derive a display name from the filename (strip extension)
					displayName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
					fontDetails = append(fontDetails, FontDetail{
						Name:     displayName,
						Files:    []string{entry.Name()},
						FlatFile: true, // flat file, served directly from fonts/
					})
				}
			}
		}
	}

	res := ResourcesResponse{
		DarkThemes:  darkThemes,
		LightThemes: lightThemes,
		Wallpapers:  wallpapers,
		Fonts:       fontDetails,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// HandleSelectWallpaper runs Zenity to select a wallpaper image natively in GNOME.
func HandleSelectWallpaper(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("zenity", "--file-selection", "--title=Select Wallpaper Image", "--file-filter=Image Files | *.png *.jpg *.jpeg *.PNG *.JPG *.JPEG")
	out, err := cmd.Output()
	if err != nil {
		// If user canceled or zenity is missing
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"canceled"}`))
		return
	}

	selectedPath := strings.TrimSpace(string(out))
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"status": "success",
		"path":   selectedPath,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// HandleWallpaperPreview serves any image file for preview in the webapp.
// Accepts ?name=<filename> for repo wallpapers, or ?path=<absolute> for custom files.
func HandleWallpaperPreview(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	path := r.URL.Query().Get("path")

	var targetPath string
	if name != "" {
		targetPath = themes.WallpaperPath(themes.Root(), name)
	} else if path != "" {
		// Security: only allow absolute paths that exist as regular files
		if !filepath.IsAbs(path) {
			http.Error(w, "Path must be absolute", http.StatusBadRequest)
			return
		}
		targetPath = path
	} else {
		http.Error(w, "name or path parameter required", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Add cache headers so repeated preview loads are fast
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, targetPath)
}

// HandleFontFile serves a font file from the themes directory.
// Supports both subdirectory layout (fonts/<Family>/<file>) and flat layout (fonts/<file>).
func HandleFontFile(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font") // display name / family name
	file := r.URL.Query().Get("file")     // actual filename
	flat := r.URL.Query().Get("flat")     // "1" if flat file layout

	if file == "" {
		http.Error(w, "file parameter is required", http.StatusBadRequest)
		return
	}

	themesRoot := themes.Root()
	fontDir := filepath.Join(themesRoot, "fonts")

	var targetPath string
	if flat == "1" || fontName == "" {
		// Flat layout: serve directly from fonts/<file>
		targetPath = filepath.Join(fontDir, file)
	} else {
		// Subdirectory layout: fonts/<Family>/<file>
		targetPath = filepath.Join(fontDir, fontName, file)
	}

	// Security: ensure path stays within themes/fonts/
	cleanFontDir := filepath.Clean(fontDir)
	cleanTarget := filepath.Clean(targetPath)
	if !strings.HasPrefix(cleanTarget, cleanFontDir+string(filepath.Separator)) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		http.Error(w, "Font file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, targetPath)
}
