package web

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"
)

// StartServer starts the HTTP server on the specified port and automatically opens the browser.
func StartServer(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	mux.HandleFunc("/api/resources", HandleResources)
	mux.HandleFunc("/api/config", HandleConfig)
	mux.HandleFunc("/api/apply", HandleApply)
	mux.HandleFunc("/api/export", HandleExport)
	mux.HandleFunc("/api/stream", HandleStream)
	mux.HandleFunc("/api/restart", HandleRestart)

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

	fmt.Printf("\nStarting web interface...\n")
	fmt.Printf("Open: %s\n\n", url)

	// Auto-open browser in background after a slight delay to allow the server to start
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	return http.Serve(listener, mux)
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
