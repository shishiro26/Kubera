#!/bin/sh
set -e

REPO="shishiro26/Kubera"

VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

FILE="kubera_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILE"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "Downloading kubera $VERSION..."
curl -fsSL "$URL" -o "$TMP/$FILE"
tar -xzf "$TMP/$FILE" -C "$TMP"

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
echo "Installing to $INSTALL_DIR..."
install -m 755 "$TMP/kubera" "$INSTALL_DIR/kubera"

echo "kubera $VERSION installed. Run: kubera"
