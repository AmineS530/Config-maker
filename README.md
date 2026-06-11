# ZoneRestore 🚀

**ZoneRestore** is a sleek, unified desktop configuration and session restoration utility designed for students at **Zone01** campus. It makes setting up or restoring your local Ubuntu environment after a machine reset completely painless.

Featuring both an interactive **command-line wizard (TUI)** built with Bubble Tea and a premium **glassmorphic local Web Dashboard** powered by Alpine JS.

---

## Key Features

- **Zsh & Prompt Customization**: Installs Oh-My-Zsh and romkatv's Powerlevel10k theme. Customize your shell prompt display name (PS1) and active Zsh command aliases dynamically.
- **Git Onboarding**: Fast-tracks Git global credential storage, username, and email setup.
- **GNOME Desktop Styling**: Configures Gnome interfaces, light/dark themes, sleep power timeout configurations, and keyboard layouts (US and French).
- **Desktop Backgrounds**: Browse available repository wallpapers with preview, or pick custom images natively using Gnome dialog selectors.
- **Docker Rootless**: Easily installs Docker inside your user directory without requiring sudo privileges.
- **Developer Fonts**: Copies custom interface and terminal monospace fonts (`MPLUS` and `MesloLGS NF`) to local storage and updates your GNOME terminal default profiles.
- **Import/Export Settings**: Export your choices to `~/.config/zonerestore/config.json` or download them as files to instantly restore your workspace next time.

---

## How to Run

Clone the repository and build the Go executable:

```zsh
# Build the binary
go build -o zonerestore cmd/main.go

# Start the interactive web interface directly
./zonerestore --web --port=8080

# Or run the interactive terminal wizard
./zonerestore --cli
```

Otherwise, simply run the executable `./zonerestore` to open the interactive menu selection screen.

---

## Feedback & Contributions

Feel free to fork the repository, adapt it to your customized workstation stack, or DM ideas on Discord to `a.sadik`.

***Happy Coding! :)***
