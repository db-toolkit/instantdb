#!/bin/bash
set -e

echo "🚀 Installing instant-db..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    *)
        echo "❌ Unsupported operating system: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest release version
echo "📦 Fetching latest release..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/db-toolkit/instant-db/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Failed to fetch latest version"
    exit 1
fi

BINARY_NAME="instant-db-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/db-toolkit/instant-db/releases/download/${LATEST_VERSION}/${BINARY_NAME}"

echo "📥 Downloading instant-db ${LATEST_VERSION} for ${OS}-${ARCH}..."
curl -L "$DOWNLOAD_URL" -o instant-db || {
    echo "❌ Download failed"
    exit 1
}

chmod +x instant-db

# Install
echo "📥 Installing to /usr/local/bin..."
if [ -w /usr/local/bin ]; then
    mv instant-db /usr/local/bin/
else
    sudo mv instant-db /usr/local/bin/
fi

echo ""
echo "✅ instant-db ${LATEST_VERSION} installed successfully!"
echo ""
echo "🎉 Get started:"
echo "   instant-db start"
echo ""
echo "📖 For more info: https://github.com/db-toolkit/instant-db"
