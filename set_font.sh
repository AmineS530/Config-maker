#!/bin/zsh

DisplayFont="MPLUS1p-Regular"
TerminalFont="MesloLGS NF Regular"

# Create the fonts directory if it doesn't exist
mkdir -p ~/.local/share/fonts

# Copy the font files to the local fonts directory
cp "fonts/$DisplayFont.ttf" ~/.local/share/fonts/
cp "fonts/$TerminalFont.ttf" ~/.local/share/fonts/

# Update the font cache
fc-cache -f -v ~/.local/share/fonts/ > /dev/null 2>&1

# Ask the user if they want to use a different display font
printf "Would you like to use a different font for display? (y/n): "
read CONFIRMATION
if [[ $CONFIRMATION == "y" || $CONFIRMATION == "Y" ]]; then
    if [ -f "fonts/$DisplayFont.ttf" ]; then
        echo "Display font is $DisplayFont"
        # Set the display font for the GNOME interface
        gsettings set org.gnome.desktop.interface font-name "'$DisplayFont 12'" 2>/dev/null
    else
        echo "Aborted. Font not found, please adjust set_font.sh"
    fi
else
    echo "Leaving display font on default"
fi

# Set the monospace font for the GNOME interface
gsettings set org.gnome.desktop.interface monospace-font-name "'$TerminalFont 12'" 2>/dev/null

# Get the default profile ID for gnome-terminal
PROFILE_ID=$(dconf list /org/gnome/terminal/legacy/profiles:/ | grep -oP '[^/]+')

# Set the terminal font for the default profile
dconf write /org/gnome/terminal/legacy/profiles:/:$PROFILE_ID/font "'$TerminalFont 12'" > /dev/null 2>&1
dconf write /org/gnome/terminal/legacy/profiles:/:$PROFILE_ID/use-system-font false /dev/null 2>&1
