package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"zonerestore/internal/cli"
	"zonerestore/internal/config"
	"zonerestore/internal/themes"
	"zonerestore/internal/web"
)

func main() {
	webFlag := flag.Bool("web", false, "Start the local web interface directly")
	cliFlag := flag.Bool("cli", false, "Start the interactive CLI wizard directly")
	portFlag := flag.Int("port", 8080, "Port for the local web interface")
	flag.Parse()

	// Initialize and prepare themes repository
	ctx := context.Background()
	_, _ = themes.EnsureThemes(ctx, os.Stdout)

	// If a flag was set, bypass the menu and run that mode directly
	if *webFlag {
		startWeb(*portFlag)
		return
	}
	if *cliFlag {
		cli.RunWizard()
		return
	}

	// Display interactive menu using Bubble Tea TUI
	reader := bufio.NewReader(os.Stdin)
	for {
		cli.ClearTerminal()
		choice := cli.RunMenu()
		switch choice {
		case 0: // CLI Wizard
			cli.ClearTerminal()
			cli.RunWizard()
			return
		case 1: // Web Interface
			cli.ClearTerminal()
			startWeb(*portFlag)
			return
		case 2: // Import Settings (View saved config.json)
			cli.ClearTerminal()
			cfg := config.LoadConfig()
			fmt.Println("\n\033[0;32m[SUCCESS] Configuration loaded from ~/.config/zonerestore/config.json:\033[0m")
			fmt.Printf("  • Install Oh-My-Zsh:       %t\n", cfg.Zsh.InstallOhMyZsh)
			fmt.Printf("  • Configure Git global:    %t (Name: %q, Email: %q)\n", cfg.Git.ConfigureGit, cfg.Git.GitName, cfg.Git.GitEmail)
			fmt.Printf("  • Apply Desktop Theme:     %t (Theme: %q, Mode: %q)\n", cfg.Theme.ApplyTheme, cfg.Theme.ThemeName, cfg.Theme.ThemeMode)
			fmt.Printf("  • Apply Background Wallpaper: %t (Image: %q)\n", cfg.Wallpaper.ApplyBackground, cfg.Wallpaper.BackgroundImage)
			fmt.Printf("  • Install Docker Rootless: %t\n", cfg.Docker.EnableDocker)
			fmt.Printf("  • Set Zsh Default Shell:   %t\n", cfg.Shell.EnableZshDefault)
			fmt.Printf("  • Pin Discord to favorites: %t\n", cfg.Dock.PinDiscord)
			fmt.Printf("  • Dual keyboard layout:    %t (Add Arabic: %t)\n", cfg.Keyboard.ConfigureKeyboard, cfg.Keyboard.AddArabic)
			fmt.Print("\nPress Enter to return to menu...")
			_, _ = reader.ReadString('\n')
		case 3: // Export Default Settings (Reset file)
			cli.ClearTerminal()
			cfg := config.DefaultConfig()
			if err := config.SaveConfig(cfg); err != nil {
				fmt.Printf("\033[0;31m[ERROR] Failed to export settings: %v\033[0m\n", err)
			} else {
				homeDir, _ := os.UserHomeDir()
				configFilePath := filepath.Join(homeDir, ".config", "zonerestore", "config.json")
				fmt.Printf("\033[0;32m[SUCCESS] Default settings exported successfully to: %s\033[0m\n", configFilePath)
			}
			fmt.Print("\nPress Enter to return to menu...")
			_, _ = reader.ReadString('\n')
		case 4, -1: // Exit / Aborted
			cli.ClearTerminal()
			fmt.Println("Exiting.")
			return
		}
	}
}

func startWeb(port int) {
	err := web.StartServer(port)
	if err != nil {
		fmt.Printf("Error starting web server: %v\n", err)
		os.Exit(1)
	}
}
