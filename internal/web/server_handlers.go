package web

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"

	"zonerestore/internal/executor"
)

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

// HandleRestart signals the server to shut down gracefully after sending a response.
// The browser receives the response, shows a closing screen, then the server exits.
func HandleRestart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"closing"}`))

	// Signal shutdown in background — gives the response time to reach the browser
	go func() {
		executor.FinishSetup(os.Stdout)
		shutdownCh <- struct{}{}
	}()
}


// getSystemThemes lists GTK themes from /usr/share/themes based on light/dark mode preference.
func getSystemThemes(mode string) []string {
	var list []string
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
			list = append(list, name)
		} else if mode == "2" && !strings.Contains(strings.ToLower(name), "dark") {
			if name == "Default" || name == "raleigh" {
				continue
			}
			list = append(list, name)
		}
	}
	return list
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
