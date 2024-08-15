#!/bin/zsh

# Define the image path
local_image="$HOME/Zone01_Desk_cfg/Background.jpeg"

# Function to open file explorer and select an image file
select_image_file() {
    FILE_PATH=$(zenity --file-selection --title="Select an Image File" --file-filter="Images | *.png *.jpg *.jpeg " 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "A file was selected"
        selected_file=$FILE_PATH
        return 0
    else
        return 1
    fi
}

# Check if the predefined image path exists
if [ -f "$local_image" ]; then
    echo "File exists at: $local_image"
    background_image=$local_image
else
    echo "File does not exist at: $local_image. Opening file explorer."
    echo "You can find some options at $HOME/Zone01_Desk_cfg/wallpapers"
    # Call the function to select an image file
    select_image_file
    if [ $? -eq 0 ]; then
        echo "Using selected file: $selected_file"
        background_image=$selected_file
    else
        echo "No file selected"
    fi
fi

# applies the background on both light and dark mode

gsettings set org.gnome.desktop.background picture-uri-dark "file://${background_image}" 2>/dev/null
gsettings set org.gnome.desktop.background picture-uri "file://${background_image}" 2>/dev/null
