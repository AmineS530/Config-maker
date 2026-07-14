# ZoneRestore Development Guide

Welcome to the ZoneRestore development documentation! This guide explains the modular architecture of the codebase, outlines what each file and folder does, and demonstrates how to extend the tool with new configuration features.

---

## 📂 Codebase File Structure

Here is the full tree layout of the repository:

```text
.
├── cmd/
│   └── main.go                  # Main entry point of the application
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── install.sh                   # Desktop wrapper script for building & running ZoneRestore
├── dev.md                       # Developer documentation (this file)
└── internal/
    ├── config/
    │   └── config.go            # Composed global UserConfig & JSON load/save storage logic
    ├── cli/                     # CLI bubbletea TUI interface context
    │   ├── menu.go              # Main interface selector menu TUI (CLI vs Web vs Load)
    │   ├── model.go             # TUI styles, color tokens, and state structs definition
    │   ├── views.go             # Step-by-step TUI rendering logic and text formats
    │   └── wizard.go            # TUI lifecycle loop, key input updates, logging console listener
    ├── web/                     # Web dashboard interface context
    │   ├── context.go           # Web-specific structures, SSE streaming writers, API models
    │   ├── config_handlers.go   # HTTP API endpoints for loading, applying and exporting JSON configs
    │   ├── resource_handlers.go # HTTP API endpoints for fetching system themes, fonts, wallpapers
    │   ├── server_handlers.go   # Main page controller, Alpine.js asset deliverer, log SSE stream router
    │   ├── server.go            # Web server bootstrapping and auto xdg-open browser opening logic
    │   └── templates.go         # Glassmorphism HTML page template containing CSS styling & Alpine wizard logic
    ├── executor/
    │   └── executor.go          # Central setup orchestrator running setting packages step-by-step
    ├── themes/
    │   └── themes.go            # git clone / git pull cache logic for remote ZoneRestoreThemes repo (Shared)
    ├── utils/
    │   ├── colors.go            # ANSI terminal color constants (Shared)
    │   ├── logger.go            # Styled terminal console output logging wrappers (Shared)
    │   └── sys.go               # Shared filesystem, file-copiers, text-appenders, and command stream runners (Shared)
    └── settings/                # Modular setting-specific packages
        ├── dock/
        │   ├── config.go        # Config representation (PinDiscord)
        │   └── executor.go      # Custom gsettings favorites app modification logic
        ├── docker/
        │   ├── config.go        # Config representation (EnableDocker)
        │   └── executor.go      # Rootless docker installation script wrapper
        ├── fonts/
        │   ├── config.go        # Config representation (ConfigureFonts, FontName)
        │   └── executor.go      # custom font TTF/OTF filesystem copy & terminal profile font appliers
        ├── git/
        │   ├── config.go        # Config representation (ConfigureGit, Name, Email)
        │   └── executor.go      # global git-config credentials setter
        ├── keyboard/
        │   ├── config.go        # Config representation (ConfigureKeyboard, AddArabic)
        │   └── executor.go      # dual US/FR standard layout with optional Arabic input source configuration
        ├── power/
        │   ├── config.go        # Config representation (ConfigurePower)
        │   └── executor.go      # GNOME settings AC sleep inactivity duration setter
        ├── shell/
        │   ├── config.go        # Config representation (EnableZshDefault)
        │   └── executor.go      # bashrc modifier to force launch shell as Zsh
        ├── theme/
        │   ├── config.go        # Config representation (ApplyTheme, Mode, Name)
        │   └── executor.go      # GTK / GNOME dark & light window decorations setter
        ├── wallpaper/
        │   ├── config.go        # Config representation (ApplyBackground, BackgroundImage)
        │   └── executor.go      # GNOME desktop wallpaper background applier
        └── zsh/
            ├── assets/          # Embedded configuration templates (.zshrc, .p10k.zsh)
            ├── config.go        # Config representation (InstallOhMyZsh)
            └── executor.go      # OhMyZsh downloader and template compiler
```

---

## 🛠 Adding a New Feature Setting

Adding a new workstation preference setup is simple. Here is a step-by-step tutorial using a mock configuration step **"Install VS Code Extensions"** (`vscode`):

### Step 1: Create a Settings Package
Create a new directory `internal/settings/vscode/` containing `config.go` and `executor.go`.

* **`internal/settings/vscode/config.go`**:
  ```go
  package vscode

  type Config struct {
      InstallExtensions bool     `json:"install_extensions"`
      ExtensionIDs      []string `json:"extension_ids"`
  }
  ```

* **`internal/settings/vscode/executor.go`**:
  ```go
  package vscode

  import (
      "io"
      "os/exec"
      "zonerestore/internal/utils"
  )

  func Apply(cfg Config, logger *utils.Logger, out io.Writer) error {
      if !cfg.InstallExtensions {
          return nil
      }

      logger.Info("Installing VS Code extensions...")
      for _, ext := range cfg.ExtensionIDs {
          logger.Info("Installing: %s", ext)
          cmd := exec.Command("code", "--install-extension", ext)
          _ = utils.RunCommandStream(cmd, out) // Streams installation live
      }
      logger.Success("VS Code extensions installation completed.")
      return nil
  }
  ```

### Step 2: Register in global Config
Open [internal/config/config.go](file:///home/asadik/Desktop/ZoneRestore/internal/config/config.go):
1. Import `zonerestore/internal/settings/vscode`.
2. Add your sub-config to `UserConfig`:
   ```diff
   type UserConfig struct {
       Zsh            zsh.Config       `json:"zsh"`
       Git            git.Config       `json:"git"`
       // ...
  +    VSCode         vscode.Config    `json:"vscode"`
   }
   ```
3. Update `DefaultConfig()` to set standard default choices:
   ```diff
   VSCode: vscode.Config{
       InstallExtensions: true,
       ExtensionIDs:      []string{"golang.go", "ms-azuretools.vscode-docker"},
   },
   ```

### Step 3: Call it in the Executor Sequence
Open [internal/executor/executor.go](file:///home/asadik/Desktop/ZoneRestore/internal/executor/executor.go):
1. Import `zonerestore/internal/settings/vscode`.
2. Invoke `vscode.Apply` inside the orchestrator flow:
   ```diff
   // 10. Default Shell Switch
   if err := shell.Apply(cfg.Shell, logger, out); err != nil {
       logger.Warning("Default shell configuration completed with warning: %v", err)
   }

  +// 11. VS Code Extensions
  +if err := vscode.Apply(cfg.VSCode, logger, out); err != nil {
  +    logger.Warning("VS Code setup completed with warning: %v", err)
  +}
   ```

### Step 4: Map to Front-ends
* **CLI (TUI)**: Open [internal/cli/wizard.go](file:///home/asadik/Desktop/ZoneRestore/internal/cli/wizard.go) or [views.go](file:///home/asadik/Desktop/ZoneRestore/internal/cli/views.go) to bind your setting to a toggle layout, summaries list, or step screen.
* **Web**:
  1. Open [internal/web/templates.go](file:///home/asadik/Desktop/ZoneRestore/internal/web/templates.go). Add the settings layout parameters to the Alpine.js state block (`cfg.vscode.install_extensions`) and map its checkbox toggle elements into Step 6.
  2. Open [internal/web/handlers.go](file:///home/asadik/Desktop/ZoneRestore/internal/web/handlers.go) if any custom asset validation is required before executing config setups.
