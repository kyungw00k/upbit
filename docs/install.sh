#!/bin/sh
# upbit CLI installer
# Usage: curl -sSL https://kyungw00k.dev/upbit/install.sh | sh
set -e

REPO="kyungw00k/upbit"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -sSf "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')
if [ -z "$VERSION" ]; then
  echo "Failed to fetch latest version"; exit 1
fi

echo "Installing upbit v$VERSION ($OS/$ARCH)..."

# Download
URL="https://github.com/$REPO/releases/download/v$VERSION/upbit_${OS}_${ARCH}.tar.gz"
if [ "$OS" = "windows" ]; then
  URL="https://github.com/$REPO/releases/download/v$VERSION/upbit_${OS}_${ARCH}.zip"
fi

TMP=$(mktemp -d)
curl -sSfL "$URL" -o "$TMP/upbit.tar.gz"
tar xzf "$TMP/upbit.tar.gz" -C "$TMP"

# Install
mkdir -p "$INSTALL_DIR"
mv "$TMP/upbit" "$INSTALL_DIR/upbit"
chmod +x "$INSTALL_DIR/upbit"
rm -rf "$TMP"

echo "Installed to $INSTALL_DIR/upbit (v$VERSION)"

# PATH check
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "Add to PATH:"
    echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc"
    echo "  source ~/.zshrc"
    ;;
esac
