package web

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"zonerestore/internal/config"
	"zonerestore/internal/themes"
)

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
	if cfg.Wallpaper.ApplyBackground && cfg.Wallpaper.BackgroundImage != "" && !filepath.IsAbs(cfg.Wallpaper.BackgroundImage) {
		cfg.Wallpaper.BackgroundImage = themes.WallpaperPath(themes.Root(), cfg.Wallpaper.BackgroundImage)
	}

	storedConfig = cfg
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))
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
