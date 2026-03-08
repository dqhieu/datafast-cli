#!/bin/sh
set -e

REPO="dqhieu/datafast-cli"
BINARY="datafast"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

LATEST=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$LATEST" ]; then
  echo "Failed to fetch latest release"
  exit 1
fi

URL="https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}_${OS}_${ARCH}.tar.gz"
echo "Downloading ${BINARY} ${LATEST} for ${OS}/${ARCH}..."

TMP=$(mktemp -d)
curl -sL "$URL" | tar xz -C "$TMP"

echo "Installing to ${INSTALL_DIR}/${BINARY}..."
sudo mv "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
sudo chmod +x "${INSTALL_DIR}/${BINARY}"
rm -rf "$TMP"

echo "Successfully installed ${BINARY} ${LATEST}"
echo "Run 'datafast auth login' to get started"
