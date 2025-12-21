#!/bin/bash
set -e

VERSION="7.2.4"

echo "🔨 Building Redis ${VERSION} for all platforms..."

BUILD_DIR="$(pwd)/build/redis"
mkdir -p "$BUILD_DIR"

# Download Redis source
echo "📦 Downloading Redis ${VERSION}..."
cd "$BUILD_DIR"
curl -L "https://download.redis.io/releases/redis-${VERSION}.tar.gz" -o redis.tar.gz
tar xzf redis.tar.gz
cd "redis-${VERSION}"

BINARIES_DIR="$BUILD_DIR/binaries"
mkdir -p "$BINARIES_DIR"

# Build for macOS (universal binary)
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "🍎 Building macOS universal binary..."
    make clean
    make -j$(sysctl -n hw.ncpu) CFLAGS="-arch arm64 -arch x86_64" LDFLAGS="-arch arm64 -arch x86_64"
    
    tar czf "$BINARIES_DIR/redis-${VERSION}-darwin-universal.tar.gz" -C src redis-server
    echo "✅ Built: redis-${VERSION}-darwin-universal.tar.gz"
    
    # Verify it's universal
    file src/redis-server
fi

# Build for Linux (requires Docker or Linux machine)
echo ""
echo "🐧 To build for Linux, run on a Linux machine or use Docker:"
echo "   docker run --rm -v \$(pwd):/work -w /work/redis-${VERSION} ubuntu:22.04 bash -c 'apt-get update && apt-get install -y build-essential && make -j\$(nproc) && tar czf /work/binaries/redis-${VERSION}-linux-amd64.tar.gz -C src redis-server'"

# Build for Windows (requires cross-compilation or Windows machine)
echo ""
echo "🪟 To build for Windows, use WSL or cross-compile with MinGW"

echo ""
echo "📍 Binaries location: $BINARIES_DIR"
echo ""
echo "📤 Upload to GitHub Releases:"
echo "   https://github.com/db-toolkit/instant-db/releases"

