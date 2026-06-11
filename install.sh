#!/bin/bash

set -euo pipefail

REPO="AmineS530/Config-maker"
BINARY_NAME="zonerestore"
INSTALL_DIR="$HOME/bin"

echo "Fetching latest release..."

VERSION=$(curl -fsSL \
  "https://api.github.com/repos/${REPO}/releases/latest" \
  | sed -n 's/.*"tag_name": "\(.*\)".*/\1/p')

echo "Latest version: $VERSION"

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

mkdir -p "$INSTALL_DIR"

echo "Downloading ${BINARY_NAME} ${VERSION}..."

curl -fL \
  -o "$INSTALL_DIR/$BINARY_NAME" \
  "$DOWNLOAD_URL"

chmod +x "$INSTALL_DIR/$BINARY_NAME"

PATH_LINE='export PATH="$HOME/bin:$PATH"'

if [ -f "$HOME/.zshrc" ] && ! grep -Fxq "$PATH_LINE" "$HOME/.zshrc"; then
    echo "$PATH_LINE" >> "$HOME/.zshrc"
fi

export PATH="$HOME/bin:$PATH"

echo
echo "Successfully installed ${BINARY_NAME} ${VERSION}"
echo

exec "$INSTALL_DIR/$BINARY_NAME"
