#!/bin/bash
set -e

echo "🚀 Installing instant-db..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go $REQUIRED_VERSION or higher is required. Found: $GO_VERSION"
    exit 1
fi

# Clone repository
echo "📦 Cloning repository..."
TEMP_DIR=$(mktemp -d)
git clone https://github.com/db-toolkit/instant-db.git "$TEMP_DIR" 2>/dev/null || {
    echo "❌ Failed to clone repository"
    exit 1
}

cd "$TEMP_DIR"

# Build
echo "🔨 Building instant-db..."
go build -o instant-db src/instantdb/cmd/instantdb/main.go || {
    echo "❌ Build failed"
    exit 1
}

# Install
echo "📥 Installing to /usr/local/bin..."
if [ -w /usr/local/bin ]; then
    mv instant-db /usr/local/bin/
else
    sudo mv instant-db /usr/local/bin/
fi

# Cleanup
cd -
rm -rf "$TEMP_DIR"

echo ""
echo "✅ instant-db installed successfully!"
echo ""
echo "🎉 Get started:"
echo "   instant-db start"
echo ""
echo "📖 For more info: https://github.com/db-toolkit/instant-db"
