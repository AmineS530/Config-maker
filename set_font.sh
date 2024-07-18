#!/bin/zsh

DisplayFont="MPLUS1p-Regular"
TerminalFont="MesloLGS NF Regular"

# Create the fonts directory if it doesn't exist
mkdir -p ~/.local/share/fonts >/dev/null 2>&1

# Copy the font files to the local fonts directory
cp "fonts/$DisplayFont.ttf" ~/.local/share/fonts/
cp "fonts/$TerminalFont.ttf" ~/.local/share/fonts/

# Update the font cache silently
fc-cache -f -v ~/.local/share/fonts/ >/dev/null 2>&1

# Ask the user if they want to use a different display font
printf "\033[0;33mWould you like to use a different font for display?\033[0m (y/n): "
read CONFIRMATION
if [[ $CONFIRMATION == "y" || $CONFIRMATION == "Y" ]]; then
    if [ -f "fonts/$DisplayFont.ttf" ]; then
        echo -e "\033[32mDisplay font is $DisplayFont\033[0m"
        # Set the display font for the GNOME interface
        gsettings set org.gnome.desktop.interface font-name "'$DisplayFont 12'" 2>/dev/null
    else
        echo -e "\033[31mAborted. Font not found, please adjust set_font.sh\033[0m"
    fi
else
    echo -e "\033[33mLeaving display font on default\033[0m"
fi

# Set the monospace font for the GNOME interface
gsettings set org.gnome.desktop.interface monospace-font-name "'$TerminalFont 12'" 2>/dev/null

# Get the default profile ID for gnome-terminal
PROFILE_ID=$(dconf list /org/gnome/terminal/legacy/profiles:/ | grep -oP '[^/]+') >/dev/null 2>&1

# Set the terminal font for the default profile
dconf write /org/gnome/terminal/legacy/profiles:/:$PROFILE_ID/font "'$TerminalFont 12'" >/dev/null 2>&1
dconf write /org/gnome/terminal/legacy/profiles:/:$PROFILE_ID/use-system-font false >/dev/null 2>&1
