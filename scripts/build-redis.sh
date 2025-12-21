#!/bin/bash
set -e

VERSION="7.2.4"
PLATFORMS=("darwin-arm64" "darwin-amd64" "linux-amd64")

echo "🔨 Building Redis binaries for all platforms..."

# Create build directory
BUILD_DIR="$(pwd)/build/redis"
mkdir -p "$BUILD_DIR"

# Download Redis source
echo "📦 Downloading Redis ${VERSION}..."
cd "$BUILD_DIR"
curl -L "https://download.redis.io/releases/redis-${VERSION}.tar.gz" -o redis.tar.gz
tar xzf redis.tar.gz
cd "redis-${VERSION}"

# Build for current platform
echo "🔨 Compiling Redis..."
make -j$(nproc 2>/dev/null || sysctl -n hw.ncpu)

# Package binaries
echo "📦 Packaging binaries..."
BINARIES_DIR="$BUILD_DIR/binaries"
mkdir -p "$BINARIES_DIR"

# Detect current platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        PLATFORM="darwin-arm64"
    else
        PLATFORM="darwin-amd64"
    fi
elif [ "$OS" = "linux" ]; then
    PLATFORM="linux-amd64"
else
    echo "❌ Unsupported platform: $OS-$ARCH"
    exit 1
fi

# Create tarball with just redis-server
echo "📦 Creating redis-${VERSION}-${PLATFORM}.tar.gz..."
tar czf "$BINARIES_DIR/redis-${VERSION}-${PLATFORM}.tar.gz" -C src redis-server

echo "✅ Built: redis-${VERSION}-${PLATFORM}.tar.gz"
echo "📍 Location: $BINARIES_DIR/redis-${VERSION}-${PLATFORM}.tar.gz"
echo ""
echo "📤 Upload this file to GitHub Releases:"
echo "   https://github.com/db-toolkit/instant-db/releases"
echo ""
echo "💡 Or test locally:"
echo "   tar xzf $BINARIES_DIR/redis-${VERSION}-${PLATFORM}.tar.gz"
echo "   ./redis-server --version"
