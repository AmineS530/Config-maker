#!/bin/bash

set -euo pipefail

REPO="AmineS530/Config-maker"
BINARY_NAME="config-maker"
INSTALL_DIR="$HOME/bin"

echo "Fetching latest release..."

VERSION=2.1

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

mkdir -p "$INSTALL_DIR"

echo "Downloading ${BINARY_NAME} ${VERSION}..."

curl -fL \
    -o "$INSTALL_DIR/$BINARY_NAME" \
    "$DOWNLOAD_URL"

chmod +x "$INSTALL_DIR/$BINARY_NAME"

PATH_LINE='export PATH="$HOME/bin:$PATH"'

if [ -f "$HOME/.bashrc" ] && ! grep -Fxq "$PATH_LINE" "$HOME/.bashrc"; then
    echo "$PATH_LINE" >> "$HOME/.bashrc"
fi

if [ -f "$HOME/.zshrc" ] && ! grep -Fxq "$PATH_LINE" "$HOME/.zshrc"; then
    echo "$PATH_LINE" >> "$HOME/.zshrc"
fi

echo
echo "Successfully installed ${BINARY_NAME} ${VERSION}"
echo

if ! command -v config-maker >/dev/null 2>&1; then
    export PATH="$HOME/bin:$PATH"
fi

exec "$INSTALL_DIR/$BINARY_NAME"
