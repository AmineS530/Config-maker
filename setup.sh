#!/bin/zsh

image_path="$HOME/zone01-config/Background.jpeg"

download_link="your_link"

# theme list ls -d /usr/share/themes/* |xargs -L 1 basename
theme_color='Yaru-viridian-dark'

# Clones the settings repo
git clone $download_link ~/Zone01_Desktop-config && cd ~/Zone01_Desktop-config 

#Move the premade p10k settings and zshrc

if [ -f ".p10k.zsh" ]; then
mv .p10k.zsh ~/.p10k.zsh
else
    echo "Powerlevel10k settings doesnt exist, skipping..."
fi

if [ -f ".zshrc" ]; then
    mv .zshrc ~/.zshrc
else
    echo "zshrc settings doesnt exist, skipping..."
fi
# clones zsh theme 
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ~/powerlevel10k

# change display and terminal font
zsh set_font.sh

# small script to set up git for first use and to rembember your info in the future
zsh git_setup.sh

# applies the background on both light and dark mode
zsh set_background.sh

# Changes theme Color
gsettings set org.gnome.desktop.interface gtk-theme $theme_color

# forward to zsh whenever termenal auto-start bash
printf "SHELL=/bin/zsh\nexec /bin/zsh -l\n" >> ~/.bashrc
