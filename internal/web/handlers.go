package web

import (
	"config-maker/internal/config"
	"config-maker/internal/executor"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	storedConfig config.UserConfig
	configMutex  sync.Mutex
)

func init() {
	storedConfig = config.LoadConfig()
}

// SSEWriter implements io.Writer to output Server-Sent Events (SSE) live stream lines.
type SSEWriter struct {
	w http.ResponseWriter
	f http.Flusher
}

// NewSSEWriter wraps a ResponseWriter and Flusher into an io.Writer.
func NewSSEWriter(w http.ResponseWriter) (*SSEWriter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("web server does not support streaming flusher")
	}
	return &SSEWriter{w: w, f: flusher}, nil
}

func (sw *SSEWriter) Write(p []byte) (n int, err error) {
	// Split incoming output into lines and send them as SSE events
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		// Escape or strip carriage returns or backspaces for beautiful web display
		trimmed = strings.ReplaceAll(trimmed, "\r", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[0m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[0;31m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[0;32m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[0;33m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[0;36m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[1;0;31m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[1;0;32m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[1;0;33m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[1;0;36m", "")
		trimmed = strings.ReplaceAll(trimmed, "\033[1m\033[92m", "")

		_, err = fmt.Fprintf(sw.w, "data: %s\n\n", trimmed)
		if err != nil {
			return 0, err
		}
	}
	sw.f.Flush()
	return len(p), nil
}

// HandleIndex serves the single-page application wizard dashboard.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(IndexTemplate))
}

// HandleResources returns available theme and wallpaper options as JSON.
func HandleResources(w http.ResponseWriter, r *http.Request) {
	darkThemes := getSystemThemes("1")
	lightThemes := getSystemThemes("2")
	wallpapers := []string{
		"976013.jpg",
		"Rin_Shima_Level_Up_Your_Web_Apps_With_Go.png",
		"wallpaper-01.png",
	}

	res := map[string]interface{}{
		"dark_themes":  darkThemes,
		"light_themes": lightThemes,
		"wallpapers":   wallpapers,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// HandleConfig returns the current stored configuration payload.
func HandleConfig(w http.ResponseWriter, r *http.Request) {
	configMutex.Lock()
	cfg := storedConfig
	configMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfg)
}

// HandleApply receives the user configuration options via JSON POST.
func HandleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	var cfg config.UserConfig
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Resolve the background image if a wallpaper index was selected
	if cfg.ApplyBackground && cfg.BackgroundImage != "" && !filepath.IsAbs(cfg.BackgroundImage) {
		homeDir, _ := os.UserHomeDir()
		cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "wallpapers", cfg.BackgroundImage)
	}

	storedConfig = cfg
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))
}

// HandleStream sets up the SSE connection and executes ApplyConfig.
func HandleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sseWriter, err := NewSSEWriter(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	configMutex.Lock()
	cfg := storedConfig
	configMutex.Unlock()

	// Run ApplyConfig with logs streaming to SSE
	exportVal := r.URL.Query().Get("export")
	export := exportVal != "false" // default to true if not specified

	err = executor.ApplyConfig(cfg, export, sseWriter)
	if err != nil {
		_, _ = fmt.Fprintf(sseWriter, "Error applying configuration: %v\n", err)
		return
	}

	_, _ = fmt.Fprintln(sseWriter, "Finished successfully!")
}

// HandleRestart shuts down terminals and launches zenity prompts in the system.
func HandleRestart(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))

	// Run in background so the web server finishes sending the response
	go func() {
		executor.FinishSetup(os.Stdout)
	}()
}

// getSystemThemes lists GTK themes from /usr/share/themes based on light/dark mode preference.
func getSystemThemes(mode string) []string {
	var themes []string
	entries, err := os.ReadDir("/usr/share/themes")
	if err != nil {
		return []string{"Yaru-dark"}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if mode == "1" && strings.Contains(strings.ToLower(name), "dark") {
			themes = append(themes, name)
		} else if mode == "2" && !strings.Contains(strings.ToLower(name), "dark") {
			if name == "Default" || name == "raleigh" {
				continue
			}
			themes = append(themes, name)
		}
	}
	return themes
}

// HandleExport receives configuration via JSON POST and exports it to the JSON file immediately.
func HandleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cfg config.UserConfig
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.SaveConfig(cfg); err != nil {
		http.Error(w, "Failed to save config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update stored config in memory too!
	configMutex.Lock()
	storedConfig = cfg
	configMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"success"}`))
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

// getWallpaperPath finds the absolute path of a wallpaper, checking both target installation and local repo folders.
func getWallpaperPath(name string) string {
	homeDir, _ := os.UserHomeDir()
	
	// Check root background
	if name == "Background.jpeg" {
		path := filepath.Join(homeDir, "Zone01_Desk_cfg", "Background.jpeg")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		wd, err := os.Getwd()
		if err == nil {
			localPath := filepath.Join(wd, "Background.jpeg")
			if _, err := os.Stat(localPath); err == nil {
				return localPath
			}
		}
		return ""
	}

	path := filepath.Join(homeDir, "Zone01_Desk_cfg", "wallpapers", name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// Fallback to local wallpapers directory
	wd, err := os.Getwd()
	if err == nil {
		localPath := filepath.Join(wd, "wallpapers", name)
		if _, err := os.Stat(localPath); err == nil {
			return localPath
		}
	}
	return ""
}

// HandleDefaultConfig returns the default configuration.
func HandleDefaultConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.DefaultConfig()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfg)
}

// HandleImportConfig opens a Zenity file dialog to select a configuration JSON file,
// parses it, updates the in-memory config, and returns it to the client.
func HandleImportConfig(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("zenity", "--file-selection", "--title=Select Configuration JSON", "--file-filter=JSON Files | *.json")
	out, err := cmd.Output()
	if err != nil {
		// User canceled or zenity is missing
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"canceled"}`))
		return
	}

	selectedPath := strings.TrimSpace(string(out))
	fileData, err := os.ReadFile(selectedPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status":"error","message":"Failed to read file"}`))
		return
	}

	var cfg config.UserConfig
	if err := json.Unmarshal(fileData, &cfg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"error","message":"Invalid JSON configuration format"}`))
		return
	}

	// Update stored config in memory
	configMutex.Lock()
	storedConfig = cfg
	configMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"status": "success",
		"config": cfg,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

//go:embed js/alpine.min.js
var alpineJS embed.FS

// HandleAlpineJS serves the embedded Alpine.js file.
func HandleAlpineJS(w http.ResponseWriter, r *http.Request) {
	data, err := alpineJS.ReadFile("js/alpine.min.js")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	_, _ = w.Write(data)
}
