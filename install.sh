#!/bin/bash

set -euo pipefail

REPO="AmineS530/ZoneRestore"
BINARY_NAME="zonerestore"
INSTALL_DIR="$HOME/bin"

echo "Fetching latest release..."

VERSION=$(curl -fsSL \
    "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | cut -d '"' -f 4)

if [ -z "$VERSION" ]; then
    echo "Failed to fetch latest release."
    exit 1
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

mkdir -p "$INSTALL_DIR"

echo "Downloading ${BINARY_NAME} ${VERSION}..."

TMP_FILE=$(mktemp)

curl -fL \
    -o "$TMP_FILE" \
    "$DOWNLOAD_URL"

chmod +x "$TMP_FILE"

mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"

PATH_LINE='export PATH="$HOME/bin:$PATH"'

if [ -f "$HOME/.zshrc" ] && ! grep -Fxq "$PATH_LINE" "$HOME/.zshrc"; then
    echo "$PATH_LINE" >> "$HOME/.zshrc"
fi

export PATH="$HOME/bin:$PATH"

echo
echo "Successfully installed ${BINARY_NAME} ${VERSION}"
echo

exec "$INSTALL_DIR/$BINARY_NAME"