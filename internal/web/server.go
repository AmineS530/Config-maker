package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"
)

// shutdownCh is signaled by HandleRestart to trigger a graceful server shutdown.
var shutdownCh = make(chan struct{}, 1)

// StartServer starts the HTTP server on the specified port and automatically opens the browser.
// It returns when the server shuts down (either from an error or from HandleRestart).
func StartServer(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	mux.HandleFunc("/api/resources", HandleResources)
	mux.HandleFunc("/api/config", HandleConfig)
	mux.HandleFunc("/api/apply", HandleApply)
	mux.HandleFunc("/api/save", HandleSave)
	mux.HandleFunc("/api/export", HandleExport) // alias for /api/save, kept for SSE executor compat
	mux.HandleFunc("/api/config/download", HandleDownload)
	mux.HandleFunc("/api/config/upload", HandleUploadConfig)
	mux.HandleFunc("/api/config/import", HandleImportConfig) // loads saved file from disk
	mux.HandleFunc("/api/config/default", HandleDefaultConfig)
	mux.HandleFunc("/api/stream", HandleStream)
	mux.HandleFunc("/api/restart", HandleRestart)
	mux.HandleFunc("/api/select-wallpaper", HandleSelectWallpaper)
	mux.HandleFunc("/api/wallpaper/preview", HandleWallpaperPreview)
	mux.HandleFunc("/api/fonts/file", HandleFontFile)
	mux.HandleFunc("/js/alpine.min.js", HandleAlpineJS)

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	// Create listener manually to handle fallback if port is already in use
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Attempt port 8081 if 8080 is taken
		addr = fmt.Sprintf("127.0.0.1:%d", port+1)
		url = fmt.Sprintf("http://%s", addr)
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to bind to port %d or %d: %w", port, port+1, err)
		}
	}

	srv := &http.Server{Handler: mux}

	fmt.Printf("\nStarting web interface...\n")
	fmt.Printf("Open: %s\n\n", url)

	// Auto-open browser in background after a slight delay to allow the server to start
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	// Watch for shutdown signal from HandleRestart
	go func() {
		<-shutdownCh
		// Give the browser a moment to receive the final HTTP response before we close
		time.Sleep(1200 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// openBrowser attempts to launch the default web browser on Linux.
func openBrowser(url string) {
	// Linux standard for opening URLs
	err := exec.Command("xdg-open", url).Start()
	if err != nil {
		// Fallback to gio open
		_ = exec.Command("gio", "open", url).Start()
	}
}
