#!/bin/bash

BOLD="$(tput bold 2>/dev/null || printf '')"
GREY="$(tput setaf 7 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
BLUE="$(tput setaf 4 2>/dev/null || printf '')"
MAGENTA="$(tput setaf 5 2>/dev/null || printf '')"
NO_COLOR="$(tput sgr0 2>/dev/null || printf '')"
LIGHT_GREEN="$(tput setaf 10 2>/dev/null || printf '')"
CYAN="$(tput setaf 37 2>/dev/null || printf '')"


info() {
  printf '%s\n' "${CYAN}==> $* ${NO_COLOR} "
}

warn() {
  printf '%s\n' "${YELLOW}! $*${NO_COLOR}"
}

error() {
  printf '%s\n' "${RED}x $*${NO_COLOR}" >&2
}

message() {
  printf '%s\n' "${GREEN}✓ $* ${NO_COLOR} "
}

sync_dotfiles(){
  message "Syncing dotfiles..."
  # Define source and destination directories
  local CONFIG_DIR="$HOME/.config"
  local DOTFILES_DIR="$HOME/dotfiles"

  # List of selected configuration directories to sync
  local configs=("fish" "kitty" "conky" "fastfetch" "i3" "qtile")

  # Sync each selected directory using rsync
  for folder in "${configs[@]}"; do
      local SRC_PATH="$CONFIG_DIR/$folder"
      local DEST_PATH="$DOTFILES_DIR/$folder"
      
      if [ -d "$SRC_PATH" ]; then
        if [ -d "$DEST_PATH" ]; then
          info "Syncing $folder..."
          # The trailing slash on the source ensures the contents are copied into the destination directory.
          rsync -av --delete "$SRC_PATH/" "$DEST_PATH/"
        else
          warn "\nDirectory $DEST_PATH does not exist, skipping syncing $folder"
        fi
      else
          warn "Directory $SRC_PATH does not exist, skipping syncing $folder..."
      fi
  done

  # sync ~/.fonts to ~/dotfiles/fonts
  local FONTS_SRC="$HOME/.fonts"
  local FONTS_DES="$HOME/dotfiles/fonts"

  mkdir -p $FONTS_DES
  mkdir -p $FONTS_SRC

  if [ ! -d "$FONTS_SRC" ]; then
    warn "Directory $FONTS_SRC does not exist, skipping fonts sync..."
  elif [ ! -d "$FONTS_DES" ]; then
    warn "Directory $FONTS_DES does not exist, skipping fonts sync..."
  else
    info "Syncing fonts..."
    rsync -av --delete "$FONTS_SRC/" "$FONTS_DES/"
  fi

  # sync ~/Pictures/wallpapers/ to ~/dotfiles/wallpapers/
  local WALLPAPER_SRC="$HOME/Pictures/wallpapers"
  local WALLPAPER_DES="$HOME/dotfiles/wallpapers"

  if [ ! -d "$WALLPAPER_SRC" ]; then
    warn "Directory $WALLPAPER_SRC does not exist, skipping wallpapers sync..."
  elif [ ! -d "$WALLPAPER_DES" ]; then
    info "Directory $WALLPAPER_DES does not exist, skipping wallpapers sync..."
  else
    info "Syncing wallpapers..."
    rsync -av --delete "$WALLPAPER_SRC/" "$WALLPAPER_DES/"
  fi

  # Sync starship.toml file from ~/.config to ~/dotfiles
  local STARSHIP_SRC="$CONFIG_DIR/starship.toml"
  local STARSHIP_DEST="$DOTFILES_DIR/starship.toml"

  if [ ! -f "$STARSHIP_SRC" ]; then
    warn "File $STARSHIP_SRC does not exist, skipping starship.toml sync..."
  elif [ ! -f "$STARSHIP_DEST" ]; then
    warn "File $STARSHIP_DEST does not exist, skipping starship.toml sync..."
  else
    info "Syncing starship.toml..."
    cp "$STARSHIP_SRC" "$STARSHIP_DEST"
  fi

  message "Sync complete."
}

sync_dotfiles

