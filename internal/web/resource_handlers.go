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

// HandleResources returns available theme, wallpaper, and font options as JSON.
func HandleResources(w http.ResponseWriter, r *http.Request) {
	darkThemes := getSystemThemes("1")
	lightThemes := getSystemThemes("2")

	themesRoot := themes.Root()
	wallpapers := themes.ListWallpapers(themesRoot)
	if len(wallpapers) == 0 {
		wallpapers = []string{
			"976013.jpg",
			"Rin_Shima_Level_Up_Your_Web_Apps_With_Go.png",
			"wallpaper-01.png",
		}
	}

	// List fonts and their files
	var fontDetails []FontDetail
	fontsList := themes.ListFonts(themesRoot)
	for _, fName := range fontsList {
		fontDir := filepath.Join(themesRoot, "fonts", fName)
		entries, err := os.ReadDir(fontDir)
		if err != nil {
			continue
		}
		var files []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".ttf" || ext == ".otf" || ext == ".woff" || ext == ".woff2" {
				files = append(files, entry.Name())
			}
		}
		if len(files) > 0 {
			fontDetails = append(fontDetails, FontDetail{Name: fName, Files: files})
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
func HandleWallpaperPreview(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	name := r.URL.Query().Get("name")

	var targetPath string
	if name != "" {
		targetPath = getWallpaperPath(name)
	} else {
		targetPath = path
	}

	if targetPath == "" {
		http.Error(w, "Path or Name parameter is required", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		http.Error(w, "Invalid file path", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, targetPath)
}

// getWallpaperPath finds the absolute path of a wallpaper in the themes repository.
func getWallpaperPath(name string) string {
	return themes.WallpaperPath(themes.Root(), name)
}

// HandleFontFile serves any font file from the themes directory.
func HandleFontFile(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	file := r.URL.Query().Get("file")
	if fontName == "" || file == "" {
		http.Error(w, "Font and file parameters are required", http.StatusBadRequest)
		return
	}
	// Safe join to prevent directory traversal
	themesRoot := themes.Root()
	targetPath := filepath.Join(themesRoot, "fonts", fontName, file)

	// Ensure the path is within the themes directory
	if !strings.HasPrefix(targetPath, filepath.Join(themesRoot, "fonts")) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		http.Error(w, "Font file not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, targetPath)
}
