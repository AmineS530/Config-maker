#!/bin/zsh
#* Made by asadik with chatGPT and love *#

# after forking and adding ur custom background / adjusting the script to your preference put the link of your repo here (your repo should be public)
download_link="https://github.com/AmineS530/Config-maker.git"

# Clones the settings repo
git clone $download_link ~/Zone01_Desk_cfg && cd ~/Zone01_Desk_cfg

#Move the premade p10k settings and zshrc

# Check for .p10k.zsh file
if [ -f ".p10k.zsh" ]; then
    mv .p10k.zsh ~/.p10k.zsh
    echo -e "\033[0;32mMoved .p10k.zsh to ~/.p10k.zsh\033[0m"
else
    echo -e "\033[0;33mPowerlevel10k settings don't exist, skipping...\033[0m"
fi

# Check for .zshrc file
if [ -f ".zshrc" ]; then
    mv .zshrc ~/.zshrc
    echo -e "\033[0;32mMoved .zshrc to ~/.zshrc\033[0m"
else
    echo -e "\033[0;33m.zshrc settings don't exist, skipping...\033[0m"
fi

# clone terminal theme
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ~/powerlevel10k
echo 'source ~/powerlevel10k/powerlevel10k.zsh-theme' >>~/.zshrc

# Function to prompt for yes/no confirmation
#!/bin/zsh

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to prompt for yes/no confirmation
confirm() {
    echo -e "${CYAN}$1${NC} ${YELLOW}(y/n):${NC} \c"
    read yn
    case $yn in
    [Yy]*) return 0 ;; # Return true if yes
    [Nn]*) return 1 ;; # Return false if no
    *)
        echo "${RED}Please answer yes or no.${NC}"
        return 1
        ;;
    esac
}

# Prompt and execute set_font.sh
if confirm "Change display and terminal fonts?\n recommended to at least change terminal font"; then
    zsh set_font.sh
fi

# Prompt and execute git_setup.sh
if confirm "Set up your Git/ea infos for first use"; then
    zsh git_setup.sh
fi

# Prompt and execute set_background.sh
if confirm "Apply A background on both light and dark modes"; then
    zsh set_background.sh
fi

# Prompt and execute set_theme.sh
if confirm "Change the GNOME theme"; then
    zsh set_theme.sh
fi

# forward to zsh whenever termenal auto-start bash
printf "SHELL=/bin/zsh\nexec /bin/zsh -l\n" >>~/.bashrc
