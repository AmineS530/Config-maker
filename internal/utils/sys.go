// SHARED UTILITY: This package is shared between the CLI interface (internal/cli)
// and the Web server (internal/web). Modifying public APIs will impact both contexts.

package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// RunCommandStream executes a command and prints its stdout/stderr live to the writer.
func RunCommandStream(cmd *exec.Cmd, out io.Writer) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Read stdout and stderr concurrently
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Fprintln(out, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Fprintln(out, scanner.Text())
		}
	}()

	return cmd.Wait()
}

// CopyFileIfExists performs file copies if source exists.
func CopyFileIfExists(src, dst string, logger *Logger) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil // skip silently if doesn't exist
	}
	return CopyFile(src, dst)
}

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// AppendToFile appends a string to the specified file.
func AppendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// GetGnomeTerminalProfiles returns active gnome terminal legacy profiles UUID strings.
func GetGnomeTerminalProfiles() ([]string, error) {
	outBytes, err := exec.Command("dconf", "list", "/org/gnome/terminal/legacy/profiles:/").Output()
	if err != nil {
		return nil, err
	}

	var profiles []string
	lines := strings.Split(string(outBytes), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		// Strip trailing slashes
		trimmed = strings.TrimSuffix(trimmed, "/")
		profiles = append(profiles, trimmed)
	}
	return profiles, nil
}

// AppendZshToBashrc appends default Zsh shell settings to .bashrc.
func AppendZshToBashrc(bashrcPath string) error {
	content, err := os.ReadFile(bashrcPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	zshLines := "SHELL=/bin/zsh\nexec /bin/zsh -l\n"
	if strings.Contains(string(content), "SHELL=/bin/zsh") {
		return nil // already exists
	}

	return AppendToFile(bashrcPath, zshLines)
}

// RemoveZshFromBashrc removes Zsh shell settings from .bashrc.
func RemoveZshFromBashrc(bashrcPath string) error {
	content, err := os.ReadFile(bashrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "SHELL=/bin/zsh" || trimmed == "exec /bin/zsh -l" {
			continue
		}
		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	if string(content) == newContent {
		return nil
	}

	return os.WriteFile(bashrcPath, []byte(newContent), 0o644)
}
