package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"zonerestore/internal/config"
	"zonerestore/internal/themes"
)

// HandleConfig returns the current stored configuration payload as JSON.
func HandleConfig(w http.ResponseWriter, r *http.Request) {
	configMutex.Lock()
	cfg := storedConfig
	configMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfg)
}

// HandleApply receives the user configuration options via JSON POST and stores them in memory.
// It does NOT persist to disk — call HandleSave or HandleDownload for that.
func HandleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cfg config.UserConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Resolve repo wallpaper filename to an absolute path
	if cfg.Wallpaper.ApplyBackground && cfg.Wallpaper.BackgroundImage != "" && !filepath.IsAbs(cfg.Wallpaper.BackgroundImage) {
		cfg.Wallpaper.BackgroundImage = themes.WallpaperPath(themes.Root(), cfg.Wallpaper.BackgroundImage)
	}

	configMutex.Lock()
	storedConfig = cfg
	configMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"success"}`))
}

// HandleSave persists the current in-memory configuration to ~/.config/zonerestore/config.json.
// Called by the web UI "Save settings" action (POST with JSON body).
func HandleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cfg config.UserConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.SaveConfig(cfg); err != nil {
		http.Error(w, "Failed to save config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	configMutex.Lock()
	storedConfig = cfg
	configMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"success"}`))
}

// HandleExport is an alias for HandleSave kept for backward compatibility with the executor stream.
func HandleExport(w http.ResponseWriter, r *http.Request) {
	HandleSave(w, r)
}

// HandleDownload serves the current configuration as a downloadable JSON file.
// This is a GET endpoint — the browser will receive a file download prompt.
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	configMutex.Lock()
	cfg := storedConfig
	configMutex.Unlock()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="zonerestore-config.json"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	_, _ = w.Write(data)
}

// HandleDefaultConfig returns the factory-default configuration as JSON.
func HandleDefaultConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.DefaultConfig()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfg)
}

// HandleUploadConfig accepts a multipart/form-data JSON file upload, parses it,
// updates the in-memory config, and returns the parsed config to the client.
// This replaces the old zenity-based HandleImportConfig.
func HandleUploadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 10 MB max (config files will be tiny, this is just a safety limit)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("config")
	if err != nil {
		http.Error(w, "No file uploaded: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	var cfg config.UserConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"error","message":"Invalid JSON configuration format"}`))
		return
	}

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

// HandleImportConfig is kept for backward compatibility but now reads from disk path
// (the saved config file) rather than launching a zenity dialog.
func HandleImportConfig(w http.ResponseWriter, r *http.Request) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Cannot determine home dir", http.StatusInternalServerError)
		return
	}
	configPath := filepath.Join(homeDir, ".config", "zonerestore", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"not_found"}`))
		return
	}

	var cfg config.UserConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"error","message":"Saved config is malformed"}`))
		return
	}

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
