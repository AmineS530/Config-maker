package main

import (
	"bufio"
	"config-maker/internal/cli"
	"config-maker/internal/config"
	"config-maker/internal/web"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	webFlag := flag.Bool("web", false, "Start the local web interface directly")
	cliFlag := flag.Bool("cli", false, "Start the interactive CLI wizard directly")
	portFlag := flag.Int("port", 8080, "Port for the local web interface")
	flag.Parse()

	// If a flag was set, bypass the menu and run that mode directly
	if *webFlag {
		startWeb(*portFlag)
		return
	}
	if *cliFlag {
		cli.RunWizard()
		return
	}

	// Display interactive menu
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nSelect interface:")
		fmt.Println("\n1) CLI Wizard")
		fmt.Println("2) Web Interface")
		fmt.Println("3) Import Settings (View saved config.json)")
		fmt.Println("4) Export Default Settings (Reset file)")
		fmt.Println("5) Exit")
		fmt.Print("\nEnter choice (1-5): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Exiting.")
			os.Exit(1)
		}

		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			cli.RunWizard()
			return
		case "2":
			startWeb(*portFlag)
			return
		case "3":
			cfg := config.LoadConfig()
			fmt.Println("\n\033[0;32m[SUCCESS] Configuration loaded from ~/.config/config-maker/config.json:\033[0m")
			fmt.Printf("  • Install Oh-My-Zsh:       %t\n", cfg.InstallOhMyZsh)
			fmt.Printf("  • Configure Git global:    %t (Name: %q, Email: %q)\n", cfg.ConfigureGit, cfg.GitName, cfg.GitEmail)
			fmt.Printf("  • Apply Desktop Theme:     %t (Theme: %q, Mode: %q)\n", cfg.ApplyTheme, cfg.ThemeName, cfg.ThemeMode)
			fmt.Printf("  • Apply Background Wallpaper: %t (Image: %q)\n", cfg.ApplyBackground, cfg.BackgroundImage)
			fmt.Printf("  • Install Docker Rootless: %t\n", cfg.EnableDocker)
			fmt.Printf("  • Set Zsh Default Shell:   %t\n", cfg.EnableZshDefault)
		case "4":
			cfg := config.DefaultConfig()
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("\033[0;31m[ERROR] Failed to find home directory: %v\033[0m\n", err)
				continue
			}
			configDir := filepath.Join(homeDir, ".config", "config-maker")
			_ = os.MkdirAll(configDir, 0755)
			configFilePath := filepath.Join(configDir, "config.json")
			configData, _ := json.MarshalIndent(cfg, "", "  ")
			if err := os.WriteFile(configFilePath, configData, 0644); err != nil {
				fmt.Printf("\033[0;31m[ERROR] Failed to export settings: %v\033[0m\n", err)
			} else {
				fmt.Printf("\033[0;32m[SUCCESS] Default settings exported successfully to: %s\033[0m\n", configFilePath)
			}
		case "5":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("\033[0;31mInvalid choice. Please select a number from 1 to 5.\033[0m")
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
