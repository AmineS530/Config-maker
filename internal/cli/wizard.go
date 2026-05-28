package cli

import (
	"bufio"
	"config-maker/internal/config"
	"config-maker/internal/executor"
	"config-maker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RunWizard kicks off the interactive command-line setup wizard.
func RunWizard() {
	reader := bufio.NewReader(os.Stdin)
	logger := utils.NewLogger()
	var cfg config.UserConfig
	homeDir, _ := os.UserHomeDir()
	configFilePath := filepath.Join(homeDir, ".config", "config-maker", "config.json")
	if _, err := os.Stat(configFilePath); err == nil {
		importSaved := promptYN(reader, "Would you like to import previously saved configuration settings?", true)
		if importSaved {
			cfg = config.LoadConfig()
			logger.Success("Imported settings successfully from config.json.")
		} else {
			cfg = config.DefaultConfig()
			logger.Info("Using default configuration settings.")
		}
	} else {
		cfg = config.DefaultConfig()
	}

	fmt.Println()
	fmt.Printf("%s==========================================%s\n", utils.Cyan, utils.Reset)
	fmt.Printf("%s      CONFIG MAKER - UBUNTU WIZARD       %s\n", utils.Cyan, utils.Reset)
	fmt.Printf("%s==========================================%s\n", utils.Cyan, utils.Reset)
	fmt.Println()

	// [Step 1/6] Oh-My-Zsh
	cfg.InstallOhMyZsh = promptYN(reader, "[Step 1/6] Install Oh-My-Zsh?", cfg.InstallOhMyZsh)

	// [Step 2/6] Configure Git
	cfg.ConfigureGit = promptYN(reader, "[Step 2/6] Configure Git?", cfg.ConfigureGit)
	if cfg.ConfigureGit {
		namePrompt := "   Enter your fullname or login:"
		if cfg.GitName != "" {
			namePrompt = fmt.Sprintf("   Enter your fullname or login [%s]:", cfg.GitName)
		}
		nameInput := promptString(reader, namePrompt)
		if nameInput != "" {
			cfg.GitName = nameInput
		}

		emailPrompt := "   Enter your email:"
		if cfg.GitEmail != "" {
			emailPrompt = fmt.Sprintf("   Enter your email [%s]:", cfg.GitEmail)
		}
		emailInput := promptString(reader, emailPrompt)
		if emailInput != "" {
			cfg.GitEmail = emailInput
		}
	}

	// [Step 3/6] Apply Theme
	cfg.ApplyTheme = promptYN(reader, "[Step 3/6] Apply Theme?", cfg.ApplyTheme)
	if cfg.ApplyTheme {
		defaultMode := "1"
		if cfg.ThemeMode == "2" {
			defaultMode = "2"
		}
		fmt.Printf("   %sDo you want Light or Dark mode?%s\n   [1] Dark\n   [2] Light\n   Enter choice (1-2) [%s]: ", utils.Yellow, utils.Reset, defaultMode)
		modeChoice := promptString(reader, "")
		if modeChoice == "" {
			modeChoice = defaultMode
		}
		if modeChoice == "2" {
			cfg.ThemeMode = "2"
		} else {
			cfg.ThemeMode = "1"
		}

		themes := getAvailableThemes(cfg.ThemeMode)
		if len(themes) > 0 {
			fmt.Printf("\n   %sAvailable themes:%s\n", utils.Cyan, utils.Reset)
			for i, t := range themes {
				fmt.Printf("   [%d] %s\n", i+1, t)
			}
			for {
				defaultThemeLabel := "1"
				if cfg.ThemeName != "" {
					defaultThemeLabel = cfg.ThemeName
				}
				fmt.Printf("   Select theme number (1-%d) [%s]: ", len(themes), defaultThemeLabel)
				idxStr := promptString(reader, "")
				if idxStr == "" {
					if cfg.ThemeName == "" {
						cfg.ThemeName = themes[0]
					}
					break
				}
				idx, err := strconv.Atoi(idxStr)
				if err == nil && idx >= 1 && idx <= len(themes) {
					cfg.ThemeName = themes[idx-1]
					break
				}
				fmt.Printf("   %sInvalid selection, try again.%s\n", utils.Red, utils.Reset)
			}
		} else {
			logger.Warning("No system themes matching selection were found. Defaulting...")
			cfg.ThemeName = "Yaru-dark"
		}
	}

	// [Step 4/6] Apply Background
	cfg.ApplyBackground = promptYN(reader, "[Step 4/6] Apply Background?", cfg.ApplyBackground)
	if cfg.ApplyBackground {
		fmt.Printf("   %sChoose wallpaper image:%s\n", utils.Yellow, utils.Reset)
		fmt.Println("   [1] Predefined Background.jpeg ($HOME/Zone01_Desk_cfg/Background.jpeg)")
		fmt.Println("   [2] Choice from repository wallpapers")
		fmt.Println("   [3] Custom absolute image path")
		fmt.Print("   Enter choice (1-3) [1]: ")
		bgChoice := promptString(reader, "")

		if bgChoice == "2" {
			wallpapers := getRepositoryWallpapers()
			if len(wallpapers) > 0 {
				fmt.Printf("\n   %sAvailable wallpapers:%s\n", utils.Cyan, utils.Reset)
				for i, wp := range wallpapers {
					fmt.Printf("   [%d] %s\n", i+1, wp)
				}
				for {
					fmt.Printf("   Select wallpaper (1-%d) [1]: ", len(wallpapers))
					idxStr := promptString(reader, "")
					var selectedWP string
					if idxStr == "" {
						selectedWP = wallpapers[0]
					} else {
						idx, err := strconv.Atoi(idxStr)
						if err == nil && idx >= 1 && idx <= len(wallpapers) {
							selectedWP = wallpapers[idx-1]
						}
					}
					if selectedWP != "" {
						homeDir, _ := os.UserHomeDir()
						cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "wallpapers", selectedWP)
						break
					}
					fmt.Printf("   %sInvalid selection, try again.%s\n", utils.Red, utils.Reset)
				}
			} else {
				fmt.Printf("   %sNo wallpapers found. Falling back to default.%s\n", utils.Yellow, utils.Reset)
				homeDir, _ := os.UserHomeDir()
				cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "Background.jpeg")
			}
		} else if bgChoice == "3" {
			for {
				path := promptString(reader, "   Enter absolute image file path:")
				if _, err := os.Stat(path); err == nil {
					cfg.BackgroundImage = path
					break
				}
				fmt.Printf("   %sFile not found at '%s'. Please enter a valid path.%s\n", utils.Red, path, utils.Reset)
			}
		} else {
			homeDir, _ := os.UserHomeDir()
			cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "Background.jpeg")
		}
	}

	// [Step 5/6] Docker
	cfg.EnableDocker = promptYN(reader, "[Step 5/6] Enable Docker Rootless?", cfg.EnableDocker)

	// [Step 6/6] Zsh Default Shell
	cfg.EnableZshDefault = promptYN(reader, "[Step 6/6] Enable Zsh as default shell in .bashrc?", cfg.EnableZshDefault)

	// Final Summary
	fmt.Println()
	fmt.Printf("%sSummary of Selections:%s\n", utils.Yellow, utils.Reset)
	printSummaryItem("Install Oh-My-Zsh", cfg.InstallOhMyZsh)
	printSummaryItem("Configure Git", cfg.ConfigureGit)
	if cfg.ConfigureGit {
		fmt.Printf("  └─ Name:  %s\n", cfg.GitName)
		fmt.Printf("  └─ Email: %s\n", cfg.GitEmail)
	}
	printSummaryItem("Apply Theme", cfg.ApplyTheme)
	if cfg.ApplyTheme {
		modeStr := "Dark"
		if cfg.ThemeMode == "2" {
			modeStr = "Light"
		}
		fmt.Printf("  └─ Mode:  %s\n", modeStr)
		fmt.Printf("  └─ Theme: %s\n", cfg.ThemeName)
	}
	printSummaryItem("Apply Background", cfg.ApplyBackground)
	if cfg.ApplyBackground {
		fmt.Printf("  └─ Image: %s\n", cfg.BackgroundImage)
	}
	printSummaryItem("Enable Docker Rootless", cfg.EnableDocker)
	printSummaryItem("Enable Zsh Default Shell", cfg.EnableZshDefault)
	fmt.Println()

	// Allow editing background at the end if settings were imported
	if _, err := os.Stat(configFilePath); err == nil && cfg.ApplyBackground {
		changeBG := promptYN(reader, "Would you like to change/edit the background image choice before applying?", false)
		if changeBG {
			fmt.Printf("   %sChoose wallpaper image:%s\n", utils.Yellow, utils.Reset)
			fmt.Println("   [1] Predefined Background.jpeg ($HOME/Zone01_Desk_cfg/Background.jpeg)")
			fmt.Println("   [2] Choice from repository wallpapers")
			fmt.Println("   [3] Custom absolute image path")
			fmt.Print("   Enter choice (1-3) [1]: ")
			bgChoice := promptString(reader, "")

			if bgChoice == "2" {
				wallpapers := getRepositoryWallpapers()
				if len(wallpapers) > 0 {
					fmt.Printf("\n   %sAvailable wallpapers:%s\n", utils.Cyan, utils.Reset)
					for i, wp := range wallpapers {
						fmt.Printf("   [%d] %s\n", i+1, wp)
					}
					for {
						fmt.Printf("   Select wallpaper (1-%d) [1]: ", len(wallpapers))
						idxStr := promptString(reader, "")
						var selectedWP string
						if idxStr == "" {
							selectedWP = wallpapers[0]
						} else {
							idx, err := strconv.Atoi(idxStr)
							if err == nil && idx >= 1 && idx <= len(wallpapers) {
								selectedWP = wallpapers[idx-1]
							}
						}
						if selectedWP != "" {
							cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "wallpapers", selectedWP)
							break
						}
						fmt.Printf("   %sInvalid selection, try again.%s\n", utils.Red, utils.Reset)
					}
				} else {
					fmt.Printf("   %sNo wallpapers found. Falling back to default.%s\n", utils.Yellow, utils.Reset)
					cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "Background.jpeg")
				}
			} else if bgChoice == "3" {
				for {
					path := promptString(reader, "   Enter absolute image file path:")
					if _, err := os.Stat(path); err == nil {
						cfg.BackgroundImage = path
						break
					}
					fmt.Printf("   %sFile not found at '%s'. Please enter a valid path.%s\n", utils.Red, path, utils.Reset)
				}
			} else {
				cfg.BackgroundImage = filepath.Join(homeDir, "Zone01_Desk_cfg", "Background.jpeg")
			}

			// Reprint the updated summary reflecting the changes
			fmt.Println()
			fmt.Printf("%sUpdated Summary of Selections:%s\n", utils.Yellow, utils.Reset)
			printSummaryItem("Install Oh-My-Zsh", cfg.InstallOhMyZsh)
			printSummaryItem("Configure Git", cfg.ConfigureGit)
			if cfg.ConfigureGit {
				fmt.Printf("  └─ Name:  %s\n", cfg.GitName)
				fmt.Printf("  └─ Email: %s\n", cfg.GitEmail)
			}
			printSummaryItem("Apply Theme", cfg.ApplyTheme)
			if cfg.ApplyTheme {
				modeStr := "Dark"
				if cfg.ThemeMode == "2" {
					modeStr = "Light"
				}
				fmt.Printf("  └─ Mode:  %s\n", modeStr)
				fmt.Printf("  └─ Theme: %s\n", cfg.ThemeName)
			}
			printSummaryItem("Apply Background", cfg.ApplyBackground)
			if cfg.ApplyBackground {
				fmt.Printf("  └─ Image: %s\n", cfg.BackgroundImage)
			}
			printSummaryItem("Enable Docker Rootless", cfg.EnableDocker)
			printSummaryItem("Enable Zsh Default Shell", cfg.EnableZshDefault)
			fmt.Println()
		}
	}

	exportConfig := promptYN(reader, "Would you like to export/save these configuration choices for future use?", true)

	apply := promptYN(reader, "Apply configuration?", true)
	if !apply {
		logger.Warning("Setup aborted. No changes made.")
		return
	}

	// Execute selections
	fmt.Println()
	err := executor.ApplyConfig(cfg, exportConfig, os.Stdout)
	if err != nil {
		logger.Error("Configuration application failed: %v", err)
		return
	}

	logger.Success("Configuration successfully applied!")

	// Reopen terminal / Finish behavior
	rebootTerminal := promptYN(reader, "Restart GNOME Terminal now to apply all changes?", true)
	if rebootTerminal {
		executor.FinishSetup(os.Stdout)
	}
}

// promptYN reads a y/n confirmation.
func promptYN(reader *bufio.Reader, label string, defaultVal bool) bool {
	choices := " (y/n) [y]: "
	if !defaultVal {
		choices = " (y/n) [n]: "
	}

	for {
		fmt.Printf("%s%s%s%s", utils.Cyan, label, utils.Reset, choices)
		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultVal
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			return defaultVal
		}
		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Printf("%sPlease answer by yes or no.%s\n", utils.Red, utils.Reset)
	}
}

// promptString prompts the user for string input.
func promptString(reader *bufio.Reader, label string) string {
	if label != "" {
		fmt.Printf("%s ", label)
	}
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// getAvailableThemes scans /usr/share/themes and filters matching names.
func getAvailableThemes(mode string) []string {
	var filtered []string
	entries, err := os.ReadDir("/usr/share/themes")
	if err != nil {
		return []string{"Yaru-dark", "Adwaita-dark"}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if mode == "1" && strings.Contains(strings.ToLower(name), "dark") {
			filtered = append(filtered, name)
		} else if mode == "2" && !strings.Contains(strings.ToLower(name), "dark") {
			// Skip basic system files
			if name == "Default" || name == "raleigh" {
				continue
			}
			filtered = append(filtered, name)
		}
	}
	return filtered
}

// getRepositoryWallpapers returns some static defaults or lists the local files.
func getRepositoryWallpapers() []string {
	// Wallpapers present in local folder or repo structure
	return []string{
		"976013.jpg",
		"Rin_Shima_Level_Up_Your_Web_Apps_With_Go.png",
		"wallpaper-01.png",
	}
}

// printSummaryItem helper to format summary lists.
func printSummaryItem(label string, enabled bool) {
	if enabled {
		fmt.Printf(" %s✔%s %s\n", utils.Green, utils.Reset, label)
	} else {
		fmt.Printf(" %s✘%s %s\n", utils.Red, utils.Reset, label)
	}
}
