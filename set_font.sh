#!/bin/zsh

DisplayFont="NimbusSanL-Reg"
TerminalFont="MesloLGS NF Regular"

# Create the fonts directory if it doesn't exist
mkdir -p ~/.local/share/fonts

# Copy the font files to the local fonts directory
cp "fonts/$DisplayFont.otf" ~/.local/share/fonts/
cp "fonts/$TerminalFont.ttf" ~/.local/share/fonts/

# Update the font cache
fc-cache -f -v ~/.local/share/fonts/

# Set the display font for the GNOME interface
gsettings set org.gnome.desktop.interface font-name "'$DisplayFont 12'"

# Set the monospace font for the GNOME interface
gsettings set org.gnome.desktop.interface monospace-font-name "'$TerminalFont 12'"

# Get the default profile ID for gnome-terminal
PROFILE_ID=$(dconf list /org/gnome/terminal/legacy/profiles:/ | grep -oP '[^/]+')

# Set the terminal font for the default profile
dconf write /org/gnome/terminal/legacy/profiles:/:$PROFILE_ID/font "'$TerminalFont 12'"
dconf write /org/gnome/terminal/legacy/profiles:/:$PROFILE_ID/use-system-font false
