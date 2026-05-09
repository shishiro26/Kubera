#!/bin/sh
# Kubera installer for Linux and macOS.
# Usage: curl -sSfL https://raw.githubusercontent.com/shishiro26/Kubera/main/install.sh | sh
set -e

REPO="shishiro26/Kubera"
BINARY="kubera"
INSTALL_DIR="$HOME/.local/bin"

# ── OS detection ──────────────────────────────────────────────────────────────
OS=$(uname -s)
case "$OS" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# ── Architecture detection ────────────────────────────────────────────────────
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)           ARCH="amd64" ;;
  aarch64 | arm64)  ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# ── Fetch latest release tag from GitHub API ──────────────────────────────────
echo "Fetching latest release..."
TAG=$(curl -sSf "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')

if [ -z "$TAG" ]; then
  echo "Could not determine latest version. Check your internet connection."
  exit 1
fi

# GoReleaser strips the leading 'v' from .Version in the archive filename.
VERSION="${TAG#v}"
ARCHIVE="kubera_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$TAG/$ARCHIVE"

# ── Download and extract ──────────────────────────────────────────────────────
echo "Downloading $ARCHIVE..."
TMP=$(mktemp -d)
curl -sSfL "$URL" | tar -xz -C "$TMP"

chmod +x "$TMP/$BINARY"

# ── Install binary ────────────────────────────────────────────────────────────
mkdir -p "$INSTALL_DIR"
mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
rm -rf "$TMP"

echo "Installed kubera $TAG to $INSTALL_DIR/$BINARY"

# ── Persist PATH in shell profile ─────────────────────────────────────────────
EXPORT_LINE="export PATH=\"$INSTALL_DIR:\$PATH\""

add_to_profile() {
  profile="$1"
  if [ -f "$profile" ] && ! grep -q "$INSTALL_DIR" "$profile" 2>/dev/null; then
    printf '\n# kubera\n%s\n' "$EXPORT_LINE" >> "$profile"
    echo "Added $INSTALL_DIR to PATH in $profile"
    return 0
  fi
  return 1
}

ADDED=0
for PROFILE in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile"; do
  if add_to_profile "$PROFILE"; then
    ADDED=1
    break
  fi
done

if [ "$ADDED" = "0" ]; then
  echo ""
  echo "Could not auto-update PATH. Add this line to your shell profile:"
  echo "  $EXPORT_LINE"
fi

echo ""
echo "Done! Open a new terminal and run: kubera init"
