#!/bin/bash
set -e

REPO="akash1551/cmdsense"
BINARY="cmdsense"
INSTALL_DIR="/usr/local/bin"

# Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

echo "Detected $OS/$ARCH..."

# Check for latest release (this API call needs a published release to work)
# LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
# For now, we'll assume a direct download link pattern or just build from source instructions if releases aren't set up.
# Since we want to support binary download:

DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/${BINARY}_${OS}_${ARCH}.tar.gz"

echo "Downloading $BINARY from $DOWNLOAD_URL..."

# Create a temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download and extract
if curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/release.tar.gz"; then
    tar -xzf "$TMP_DIR/release.tar.gz" -C "$TMP_DIR"
else
    echo "Failed to download release. Please check if '$REPO' has releases."
    exit 1
fi

# Install
echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
else
    sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
fi

echo "✅ Installed $BINARY successfully!"
echo "Run '$BINARY config' to get started."
