#!/usr/bin/env bash
set -euo pipefail

REPO="TheMightyDman/tache"
BIN="tache"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "Unsupported arch: $ARCH" >&2; exit 1 ;;
esac

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep -E '"tag_name":' | cut -d '"' -f4)
if [ -z "$LATEST" ]; then
  echo "Could not determine latest release. Ensure a release exists." >&2
  exit 1
fi

TARBALL="${BIN}_${LATEST}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT
cd "$TMPDIR"
echo "Downloading $URL" >&2
curl -fsSL -o "$TARBALL" "$URL"
tar -xzf "$TARBALL"

mkdir -p "$INSTALL_DIR"
install -m 0755 "$BIN" "$INSTALL_DIR/$BIN"

if ! command -v tache >/dev/null 2>&1; then
  echo "Installed to $INSTALL_DIR. Add it to your PATH, for example:" >&2
  echo "  export PATH=\"$INSTALL_DIR:\$PATH\"" >&2
else
  echo "Installed $BIN to $(command -v tache)" >&2
fi

