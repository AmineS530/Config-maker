#!/bin/zsh

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

#*  -------------------------------  Made by asadik with chatGPT and love  -------------------------------  *#

# after forking and adding ur custom background / adjusting the script to your preference put the link of your repo here (your repo should be public)
download_link="https://github.com/AmineS530/Config-maker.git"

# Clones the settings repo
destination_dir="$HOME/Zone01_Desk_cfg"

# Check if the directory already exists
if [ -d "$destination_dir" ]; then
    echo -e "${YELLOW}Directory $destination_dir already exists. Overwriting...${NC}"
    # Remove existing directory and clone the repository
    rm -rf "$destination_dir"
    git clone --quiet --depth=1 "$download_link" "$destination_dir" && cd "$destination_dir"
else
    echo -e "${CYAN}Cloning repository into $destination_dir...${NC}"
    # Directory does not exist, clone the repository
    git clone --quiet --depth=1 "$download_link" "$destination_dir" && cd "$destination_dir"
fi

# Move the premade p10k settings and zshrc

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
    echo 'source ~/powerlevel10k/powerlevel10k.zsh-theme' >> ~/.zshrc
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

# Prompt and execute set_font.sh
if confirm "Change display and terminal fonts?\nRecommended to at least change terminal font if using p10k"; then
    zsh set_font.sh
fi

# Prompt and execute git_setup.sh
if confirm "Set up your Git/ea configs for first use"; then
    zsh git_setup.sh
fi

# Prompt and execute set_background.sh
if confirm "Apply a background on both light and dark modes"; then
    zsh set_background.sh
fi

# Prompt and execute set_theme.sh
if confirm "Change the GNOME theme"; then
    zsh set_theme.sh
fi

# Forward to zsh whenever terminal auto-starts bash
printf "SHELL=/bin/zsh\nexec /bin/zsh -l\n" >> ~/.bashrc
