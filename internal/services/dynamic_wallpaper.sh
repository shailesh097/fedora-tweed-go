#!/bin/bash

# Change the wallpaper after a certain interval from a directory
change_wallpaper(){
    WALLPAPER_DIR="$HOME/Pictures/wallpapers"
    INTERVAL=3600
    
    # Get the list of wallpapers
    mapfile -t WALLPAPERS < <(find "$WALLPAPER_DIR" -type f | sort)
    
    # Select a random starting wallpaper
    RANDOM_INDEX=$(shuf -i 0-$((${#WALLPAPERS[@]} - 1)) -n 1)
    
    # Loop through wallpapers sequentially after the first random selection
    while true; do
        # Get the wallpaper from the list
        WALLPAPER="${WALLPAPERS[$RANDOM_INDEX]}"
        
        # Set the wallpaper using gsettings
        gsettings set org.gnome.desktop.background picture-uri-dark "file://$WALLPAPER"
        
        # Increment index and loop back to the beginning if needed
        RANDOM_INDEX=$(( (RANDOM_INDEX + 1) % ${#WALLPAPERS[@]} ))
        
        # Wait for the specified interval before changing the wallpaper again
        sleep $INTERVAL
    done
}

change_wallpaper

