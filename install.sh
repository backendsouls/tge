#!/usr/bin/env bash
set -e

# Repository and binary details
REPO="backendsouls/tge"
BINARY_NAME="tge"

echo "=== The Great Emulator (TGE) Installer ==="

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="x86_64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected OS: $OS | Architecture: $ARCH"
echo "Fetching the latest release info from GitHub..."

# Fetch the latest release JSON from GitHub API
RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")

# Extract the browser_download_url for the matching OS and ARCH
DOWNLOAD_URL=$(echo "$RELEASE_JSON" | grep -i "browser_download_url" | grep -i "$OS" | grep -i "$ARCH" | cut -d '"' -f 4 | head -n 1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Error: Could not find a pre-built binary for $OS $ARCH."
    echo "Please check the releases page manually: https://github.com/$REPO/releases"
    exit 1
fi

echo "Downloading TGE from $DOWNLOAD_URL..."
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"
curl -sL "$DOWNLOAD_URL" -o tge.tar.gz

echo "Extracting..."
tar -xzf tge.tar.gz "$BINARY_NAME"

INSTALL_DIR="/usr/local/bin"
echo "Installing to $INSTALL_DIR (this may prompt for your password via sudo)..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

cd - > /dev/null
rm -rf "$TMP_DIR"

echo "✅ Success! TGE has been successfully installed."
echo "Run 'tge --help' to get started."
