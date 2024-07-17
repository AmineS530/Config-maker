#!/bin/zsh

# Define color variables
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# List of available themes
themes=("Adwaita" "Adwaita-dark" "Default" "Emacs" "HighContrast" "Yaru" "Yaru-bark" "Yaru-bark-dark" "Yaru-blue" "Yaru-blue-dark" "Yaru-dark" "Yaru-dark-hdpi" "Yaru-dark-xhdpi" "Yaru-hdpi" "Yaru-magenta" "Yaru-magenta-dark" "Yaru-olive" "Yaru-olive-dark" "Yaru-prussiangreen" "Yaru-prussiangreen-dark" "Yaru-purple" "Yaru-purple-dark" "Yaru-red" "Yaru-red-dark" "Yaru-sage" "Yaru-sage-dark" "Yaru-viridian" "Yaru-viridian-dark" "Yaru-xhdpi")

# Prompt the user to select a theme
echo -e "${BLUE}Available themes:${NC}"
for i in {2..${#themes[@]}}; do
    echo -e "${GREEN}[$((i - 1))] ${themes[$((i - 1))]}${NC}"
done

printf "${YELLOW}Enter the number corresponding to the theme you want to apply: ${NC}"
read theme_index

# Validate input
if [[ "$theme_index" =~ ^[0-9]+$ && $theme_index -ge 0 && $theme_index -lt ${#themes[@]} ]]; then
    selected_theme=${themes[$theme_index]}
    echo -e "${GREEN}Applying theme: '$selected_theme'${NC}"

    # Apply the selected theme
    gsettings set org.gnome.desktop.interface gtk-theme "'$selected_theme'" 2>/dev/null
    gsettings set org.gnome.desktop.wm.preferences theme "'$selected_theme'" 2>/dev/null
    sleep 1

    # Optionally set a dark color scheme if the selected theme contains "dark"
    if [[ "${selected_theme}" =~ "dark" ]]; then
        gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark' 2>/dev/null
    else
        gsettings set org.gnome.desktop.interface color-scheme 'default' 2>/dev/null
    fi
else
    echo -e "${RED}Invalid selection. Please run the script again and select a valid number.${NC}"
fi
