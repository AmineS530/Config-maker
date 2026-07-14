package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zonerestore/internal/utils"
)

func (m tuiModel) View() string {
	var s strings.Builder

	// Top Title and ASCII brand banner
	s.WriteString(activeItemStyle.Render(asciiBanner) + "\n")
	s.WriteString(titleStyle.Render("   ZONE RESTORE - CONFIGURATION WIZARD   ") + "\n\n")

	var stepContent string

	// Determine dynamic responsive box width
	boxW := m.width - 6
	if boxW < 60 {
		boxW = 60
	} else if boxW > 85 {
		boxW = 85
	}
	responsiveBoxStyle := boxStyle.Width(boxW)

	// Step Progress Bar
	if m.step > 0 && m.step < 8 {
		s.WriteString(renderProgressBar(m.step, 7) + "\n\n")
	}

	switch m.step {
	case 0: // Import Settings Prompt
		stepContent = fmt.Sprintf(
			" Loaded existing configurations in ~/.config/zonerestore/config.json.\n\n"+
				" %sWould you like to import your saved settings?%s\n\n"+
				"   %s\n   %s",
			utils.Yellow, utils.Reset,
			renderToggleOption("Yes (Import Settings)", m.importPrompt),
			renderToggleOption("No (Pristine Defaults)", !m.importPrompt),
		)

	case 1: // Oh-My-Zsh Setup
		stepContent = fmt.Sprintf(
			" [Step 1/7] Oh-My-Zsh Shell Framework\n\n"+
				" %sWould you like to install Oh-My-Zsh unattended?%s\n\n"+
				"   %s\n   %s",
			utils.Yellow, utils.Reset,
			renderToggleOption("Yes, install Oh-My-Zsh", m.ohMyZshChoice),
			renderToggleOption("No, skip shell framework installation", !m.ohMyZshChoice),
		)

	case 2: // Git Setup
		if !m.configureGit {
			stepContent = fmt.Sprintf(
				" [Step 2/7] Git Global Credentials\n\n"+
					" %sWould you like to set up Git global name and email?%s\n\n"+
					"   %s\n   %s",
				utils.Yellow, utils.Reset,
				renderToggleOption("Yes, set up Git credentials", m.configureGit),
				renderToggleOption("No, skip Git setup", !m.configureGit),
			)
		} else {
			nameLabel := "   Full Name:     "
			emailLabel := "   Email Address: "
			if m.focusedInput == 0 {
				nameLabel = " " + activeItemStyle.Render("▶ Full Name:     ")
			} else {
				emailLabel = " " + activeItemStyle.Render("▶ Email Address: ")
			}

			stepContent = fmt.Sprintf(
				" [Step 2/7] Git Global Credentials\n\n"+
					" %sPlease specify your git identity:%s\n\n"+
					"%s%s\n"+
					"%s%s\n\n"+
					" %s",
				utils.Yellow, utils.Reset,
				nameLabel, m.gitNameInput.View(),
				emailLabel, m.gitEmailInput.View(),
				helpStyle.Render("Use Up/Down or Tab to switch inputs. Enter to validate and continue."),
			)
		}

	case 3: // Themes Setup
		if !m.applyTheme {
			stepContent = fmt.Sprintf(
				" [Step 3/7] Gnome Desktop Theme\n\n"+
					" %sWould you like to apply GTK window decoration themes?%s\n\n"+
					"   %s\n   %s",
				utils.Yellow, utils.Reset,
				renderToggleOption("Yes, apply Gnome themes", m.applyTheme),
				renderToggleOption("No, leave desktop theme untouched", !m.applyTheme),
			)
		} else {
			modeText := renderHorizontalOptions([]string{"Dark Mode", "Light Mode"}, m.themeModeCursor)

			var list string
			if len(m.availableThemes) == 0 {
				list = fmt.Sprintf("     %sNo matching system themes detected%s\n", utils.Red, utils.Reset)
			} else {
				var items []string
				start := m.themeCursor - 2
				if start < 0 {
					start = 0
				}
				end := start + 5
				if end > len(m.availableThemes) {
					end = len(m.availableThemes)
				}
				for i := start; i < end; i++ {
					indicator := "  "
					if i == m.themeCursor {
						indicator = "▶ "
					}
					name := m.availableThemes[i]
					if i == m.themeCursor {
						items = append(items, indicator+activeItemStyle.Render(name))
					} else {
						items = append(items, indicator+name)
					}
				}
				list = "     " + strings.Join(items, "\n     ") + "\n"
			}

			stepContent = fmt.Sprintf(
				" [Step 3/7] Gnome Desktop Theme\n\n"+
					"   Color mode:  %s\n\n"+
					"   Theme list:\n%s\n"+
					"   %s",
				modeText,
				list,
				helpStyle.Render("Left/Right switches Color mode, Up/Down selects theme. Enter to validate."),
			)
		}

	case 4: // Font Selection Setup
		if !m.applyFonts {
			stepContent = fmt.Sprintf(
				" [Step 4/7] Fonts\n\n"+
					"   %s⚠  Terminal font is always: MesloLGS NF Regular 12 (locked)%s\n\n"+
					" %sWould you like to install fonts and set a display font?%s\n\n"+
					"   %s\n   %s",
				utils.Cyan, utils.Reset,
				utils.Yellow, utils.Reset,
				renderToggleOption("Yes, install & configure fonts", m.applyFonts),
				renderToggleOption("No, skip font setup", !m.applyFonts),
			)
		} else {
			var list string
			if len(m.availableFonts) == 0 {
				list = fmt.Sprintf("     %sNo custom fonts found in repo folder%s\n", utils.Red, utils.Reset)
			} else {
				var items []string
				for i, f := range m.availableFonts {
					indicator := "  "
					if i == m.fontCursor {
						indicator = "▶ "
					}
					if i == m.fontCursor {
						items = append(items, indicator+activeItemStyle.Render(f))
					} else {
						items = append(items, indicator+f)
					}
				}
				list = "     " + strings.Join(items, "\n     ") + "\n"
			}

			stepContent = fmt.Sprintf(
				" [Step 4/7] Fonts\n\n"+
					"   %s⚠  Terminal font (locked): MesloLGS NF Regular 12%s\n\n"+
					"   Select a display/UI font for GNOME (optional):\n%s\n"+
					"   [System Default = Ubuntu 11]\n\n"+
					"   %s",
				utils.Cyan, utils.Reset,
				list,
				helpStyle.Render("Up/Down to navigate. Enter to choose and proceed."),
			)
		}

	case 5: // Background Setup
		if !m.applyBg {
			stepContent = fmt.Sprintf(
				" [Step 5/7] Desktop Wallpaper\n\n"+
					" %sWould you like to apply a desktop background?%s\n\n"+
					"   %s\n   %s",
				utils.Yellow, utils.Reset,
				renderToggleOption("Yes, apply desktop wallpaper", m.applyBg),
				renderToggleOption("No, skip wallpaper setup", !m.applyBg),
			)
		} else {
			options := []string{
				"Predefined wallpaper (Background.jpeg)",
				"Choose from repository themes list",
				"Specify custom absolute image path",
			}
			var list strings.Builder
			for idx, opt := range options {
				pointer := "  "
				if idx == m.bgChoiceCursor {
					pointer = "▶ "
				}
				if idx == m.bgChoiceCursor {
					list.WriteString(pointer + activeItemStyle.Render(opt) + "\n")
				} else {
					list.WriteString(pointer + opt + "\n")
				}

				// Render list item detail expansions
				if idx == 1 && m.bgChoiceCursor == 1 {
					var wpItems []string
					for i, wp := range m.repoWallpapers {
						if i == m.wallpaperCursor {
							wpItems = append(wpItems, activeItemStyle.Render("["+wp+"]"))
						} else {
							wpItems = append(wpItems, wp)
						}
					}
					list.WriteString("     Wallpaper: " + strings.Join(wpItems, "   ") + "\n")
				}
			}

			pathInputView := ""
			if m.bgChoiceCursor == 2 {
				pathInputView = "\n   Input Path: " + m.customBgInput.View() + "\n"
			}

			stepContent = fmt.Sprintf(
				" [Step 5/7] Desktop Wallpaper\n\n"+
					"%s"+
					"%s\n"+
					"   %s",
				list.String(),
				pathInputView,
				helpStyle.Render("Press Up/Down to choose option, Left/Right to scroll wallpapers, Enter to select."),
			)
		}

	case 6: // Docker, Zsh, Keyboard layouts, Power, Dock checkboxes
		stepContent = fmt.Sprintf(
			" [Step 6/7] System preferences & defaults\n\n"+
				"   %s Enable Docker Rootless?%s\n"+
				"     %s\n\n"+
				"   %s Set Zsh as Default Shell?%s\n"+
				"     %s\n\n"+
				"   %s Pin Discord to favourites in Dock?%s\n"+
				"     %s\n\n"+
				"   %s Configure Gnome keyboard layout (US+FR)?%s\n"+
				"     %s\n\n"+
				"   %s Add Arabic layout to keyboard list?%s\n"+
				"     %s\n\n"+
				"   %s Configure standard power sleep settings?%s\n"+
				"     %s\n\n"+
				"   %s",
			renderActiveLabel("[Option A]", m.focusedInput == 0), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.enableDocker)),
			renderActiveLabel("[Option B]", m.focusedInput == 1), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.enableZshDefault)),
			renderActiveLabel("[Option C]", m.focusedInput == 2), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.pinDiscord)),
			renderActiveLabel("[Option D]", m.focusedInput == 3), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.configureKeyboard)),
			renderActiveLabel("[Option E]", m.focusedInput == 4), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.addArabic)),
			renderActiveLabel("[Option F]", m.focusedInput == 5), utils.Reset,
			renderHorizontalOptions([]string{"Yes", "No"}, getYesNoIndex(!m.configurePower)),
			helpStyle.Render("Press Up/Down/Tab to switch, Left/Right/Space to toggle, Enter to confirm."),
		)

	case 7: // Final Summary Checklist & Save
		var summary strings.Builder
		summary.WriteString(fmt.Sprintf(" %sReview Selections:%s\n\n", utils.Yellow, utils.Reset))
		summary.WriteString(renderSummaryRow("Install Oh-My-Zsh", m.ohMyZshChoice) + "\n")
		summary.WriteString(renderSummaryRow("Configure Git", m.configureGit) + "\n")
		if m.configureGit {
			summary.WriteString(fmt.Sprintf("   └─ Name:  %s\n", m.gitNameInput.Value()))
			summary.WriteString(fmt.Sprintf("   └─ Email: %s\n", m.gitEmailInput.Value()))
		}
		summary.WriteString(renderSummaryRow("Apply Gnome Theme", m.applyTheme) + "\n")
		if m.applyTheme {
			mode := "Dark"
			if m.themeModeCursor == 1 {
				mode = "Light"
			}
			theme := "Yaru-dark"
			if len(m.availableThemes) > 0 {
				theme = m.availableThemes[m.themeCursor]
			}
			summary.WriteString(fmt.Sprintf("   └─ Theme: %s (%s)\n", theme, mode))
		}
		summary.WriteString(renderSummaryRow("Install & Configure Fonts", m.applyFonts) + "\n")
		if m.applyFonts {
			summary.WriteString(fmt.Sprintf("   └─ Terminal font: MesloLGS NF Regular 12 (locked)\n"))
			displayFont := "System Default (Ubuntu 11)"
			if len(m.availableFonts) > 0 && m.fontCursor > 0 {
				displayFont = m.availableFonts[m.fontCursor]
			}
			summary.WriteString(fmt.Sprintf("   └─ Display font: %s\n", displayFont))
		}
		summary.WriteString(renderSummaryRow("Apply Wallpaper", m.applyBg) + "\n")
		if m.applyBg {
			var bg string
			if m.bgChoiceCursor == 0 {
				bg = "Predefined Background.jpeg"
			} else if m.bgChoiceCursor == 1 {
				bg = m.repoWallpapers[m.wallpaperCursor]
			} else {
				bg = m.customBgInput.Value()
			}
			summary.WriteString(fmt.Sprintf("   └─ Path:  %s\n", bg))
		}
		summary.WriteString(renderSummaryRow("Install Docker Rootless", m.enableDocker) + "\n")
		summary.WriteString(renderSummaryRow("Set Zsh Default Shell", m.enableZshDefault) + "\n")
		summary.WriteString(renderSummaryRow("Pin Discord to Dock", m.pinDiscord) + "\n")
		summary.WriteString(renderSummaryRow("Configure Keyboard layouts", m.configureKeyboard) + "\n")
		if m.configureKeyboard {
			summary.WriteString(fmt.Sprintf("   └─ Layouts: US + FR%s\n", map[bool]string{true: " + AR (Arabic)", false: ""}[m.addArabic]))
		}
		summary.WriteString(renderSummaryRow("Configure Sleep settings", m.configurePower) + "\n\n")

		summary.WriteString(helpStyle.Render("Press Enter to execute the configuration updates, or Esc/q to quit."))
		stepContent = summary.String()

	case 8: // Live Execution Console
		var logs strings.Builder
		logs.WriteString(fmt.Sprintf(" %sExecution Console:%s\n\n", utils.Yellow, utils.Reset))

		if len(m.logLines) == 0 {
			logs.WriteString("   Initialising setup environment...\n")
		} else {
			start := len(m.logLines) - 10
			if start < 0 {
				start = 0
			}
			for i := start; i < len(m.logLines); i++ {
				logs.WriteString("   " + m.logLines[i] + "\n")
			}
		}

		logs.WriteString("\n")

		if !m.finished {
			logs.WriteString(fmt.Sprintf("   %s Applying configuration selections live... Please wait.%s\n", m.spinner.View(), utils.Reset))
		} else {
			if m.err != nil {
				logs.WriteString(fmt.Sprintf("   %s✖ Failed: %v%s\n\n", utils.Red, m.err, utils.Reset))
				logs.WriteString("   Press [q] or [Esc] to exit.")
			} else {
				if !m.exportDone {
					logs.WriteString(fmt.Sprintf("   %s✔ Finished successfully!%s\n\n", utils.Green, utils.Reset))
					logs.WriteString(fmt.Sprintf("   %sWould you like to export/save these settings for future use?%s\n", utils.Yellow, utils.Reset))
					logs.WriteString(fmt.Sprintf("     %s\n     %s\n\n",
						renderToggleOption("Yes, save choices", m.exportPrompt),
						renderToggleOption("No, exit program", !m.exportPrompt),
					))
					logs.WriteString(helpStyle.Render("Use Up/Down to navigate choice. Enter to confirm selection."))
				} else {
					logs.WriteString(fmt.Sprintf("   %s✔ Settings exported. Exiting terminal.%s\n", utils.Green, utils.Reset))
				}
			}
		}
		stepContent = logs.String()
	}

	s.WriteString(responsiveBoxStyle.Render(stepContent))
	return s.String()
}

// ── RENDER HELPERS ──

func renderProgressBar(current, total int) string {
	pct := float64(current) / float64(total)
	width := 32
	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf(" %s Progress: [%s] %d%%", utils.Cyan, bar, int(pct*100))
}

func renderToggleOption(label string, active bool) string {
	if active {
		return fmt.Sprintf("[%sX%s] %s", utils.Cyan, utils.Reset, activeItemStyle.Render(label))
	}
	return fmt.Sprintf("[ ] %s", label)
}

func renderActiveLabel(label string, active bool) string {
	if active {
		return activeItemStyle.Render(label)
	}
	return label
}

func renderHorizontalOptions(options []string, selectedIndex int) string {
	var out []string
	for idx, opt := range options {
		if idx == selectedIndex {
			out = append(out, fmt.Sprintf("[%s%s%s]", utils.Cyan, opt, utils.Reset))
		} else {
			out = append(out, fmt.Sprintf(" %s ", opt))
		}
	}
	return strings.Join(out, "   ")
}

func renderSummaryRow(label string, enabled bool) string {
	status := errorStyle.Render("Disabled")
	if enabled {
		status = successStyle.Render("Enabled")
	}
	return fmt.Sprintf("  • %-32s -> %s", label, status)
}

func getYesNoIndex(isNo bool) int {
	if isNo {
		return 1
	}
	return 0
}

func getSystemThemes(mode string) []string {
	// Look inside system path for GTK themes
	var paths []string
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".themes"))
	}
	paths = append(paths, "/usr/share/themes")

	var themesList []string
	for _, p := range paths {
		dirs, err := os.ReadDir(p)
		if err != nil {
			continue
		}
		for _, d := range dirs {
			if !d.IsDir() {
				continue
			}
			name := d.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			// basic filtering by Mode name prefix
			isDark := strings.Contains(strings.ToLower(name), "dark") ||
				strings.Contains(strings.ToLower(name), "black") ||
				strings.Contains(strings.ToLower(name), "night")

			if mode == "1" && isDark {
				themesList = append(themesList, name)
			} else if mode == "2" && !isDark {
				themesList = append(themesList, name)
			}
		}
	}

	// Deduplicate
	dedup := make(map[string]bool)
	var finalThemes []string
	for _, t := range themesList {
		if !dedup[t] {
			dedup[t] = true
			finalThemes = append(finalThemes, t)
		}
	}

	if len(finalThemes) == 0 {
		if mode == "1" {
			return []string{"Yaru-dark", "Adwaita-dark"}
		}
		return []string{"Yaru", "Adwaita"}
	}
	return finalThemes
}
