#!/bin/zsh

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

#*  <-------------------------------  Made by asadik with chatGPT and love  ------------------------------->  *#

# Move the premade p10k settings and zshrc
mv .p10k.zsh ~/.p10k.zsh && mv .zshrc ~/.zshrc
# Check for .p10k.zsh file
if [ -f ".p10k.zsh" ]; then
    mv .p10k.zsh ~/.p10k.zsh
    echo -e "${GREEN}Moved .p10k.zsh to ~/.p10k.zsh${NC}"
else
    echo -e "${YELLOW}Powerlevel10k settings don't exist, skipping...${NC}"
fi

# Check for .zshrc file
if [ -f ".zshrc" ]; then
    mv .zshrc ~/.zshrc
    echo -e "${GREEN}Moved .zshrc to ~/.zshrc${NC}"
else
    echo -e "${YELLOW}.zshrc file don't exist, skipping...${NC}"
fi

# Check if the directory exists
if [ -d "$HOME/powerlevel10k" ]; then
    echo -e "${YELLOW}Powerlevel10k directory already exists, skipping clone.${NC}"
else
    # Clone terminal theme
    echo -e "${CYAN}Cloning Powerlevel10k repository...${NC}"
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ~/powerlevel10k
    echo 'source ~/powerlevel10k/powerlevel10k.zsh-theme' >>~/.zshrc
fi

# Function to prompt for yes/no confirmation
confirm() {
    echo -e "${CYAN}$1${NC} ${YELLOW}(y/n):${NC} \c"
    read yn
    case $yn in
    [Yy]*) return 0 ;; # Return true if yes
    [Nn]*) return 1 ;; # Return false if no
    *)
        echo -e "${RED}Please answer by yes or no.${NC}"
        return 1
        ;;
    esac
}

# Makes you use french and english layouts
gsettings set org.gnome.desktop.input-sources sources "[('xkb', 'us'), ('xkb', 'fr')]" > /dev/null 2>&1

# Prompt and execute set_background.sh
if confirm "Would you like to change the background?"; then
    zsh set_background.sh
fi

# Prompt and execute set_theme.sh
if confirm "Change the device theme"; then
    zsh set_theme.sh
fi

# Get git/ea configs
if confirm "setup github/gitea settings for first use"; then
    zsh git_setup.sh
fi

# Prompt and execute set_font.sh
zsh set_font.sh

# Forward to zsh whenever terminal auto-starts bash
printf "SHELL=/bin/zsh\nexec /bin/zsh -l\n" >>~/.bashrc
