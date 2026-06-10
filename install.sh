#!/bin/bash
set -e

VERSION="v0.1"
BINARY_NAME="config-maker"
INSTALL_DIR="$HOME/bin"

DOWNLOAD_URL="https://github.com/AmineS530/Config-maker/releases/download/${VERSION}/${BINARY_NAME}"

mkdir -p "$INSTALL_DIR"

echo "Downloading ${BINARY_NAME} ${VERSION}..."
curl -L -o "$INSTALL_DIR/$BINARY_NAME" "$DOWNLOAD_URL"

chmod +x "$INSTALL_DIR/$BINARY_NAME"

if ! grep -q 'export PATH="$HOME/bin:$PATH"' "$HOME/.bashrc" 2>/dev/null; then
    echo 'export PATH="$HOME/bin:$PATH"' >> "$HOME/.bashrc"
fi

echo "Installed ${BINARY_NAME} ${VERSION}"
"$INSTALL_DIR/$BINARY_NAME"
