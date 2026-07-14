package web

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"zonerestore/internal/config"
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

// FontDetail holds details about a system/downloaded developer monospace font family.
type FontDetail struct {
	Name     string   `json:"name"`
	Files    []string `json:"files"`
	FlatFile bool     `json:"flat_file"` // true when font files live directly in themes/fonts/ (not in a subdirectory)
}

// ResourcesResponse holds response metadata returned by HandleResources.
type ResourcesResponse struct {
	DarkThemes  []string     `json:"dark_themes"`
	LightThemes []string     `json:"light_themes"`
	Wallpapers  []string     `json:"wallpapers"`
	Fonts       []FontDetail `json:"fonts"`
}
