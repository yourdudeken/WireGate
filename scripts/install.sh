#!/bin/sh
set -e

# WireGate Installer
# This script installs the wiregate binary to /usr/local/bin

echo "WireGate Installer"
echo "--------------------"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${OS}" in
    linux*)  OS='linux';;
    darwin*) OS='darwin';;
    *)       echo "Error: OS ${OS} is not supported."; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64) ARCH='amd64';;
    arm64|aarch64) ARCH='arm64';;
    *)      echo "Error: Architecture ${ARCH} is not supported."; exit 1;;
esac

BINARY_NAME="wiregate"
REPO="yourdudeken/wiregate"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "Downloading ${BINARY_NAME} for ${OS}/${ARCH}..."
if ! curl -sSfL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${BINARY_NAME}"; then
    echo "Error: Failed to download binary. Please ensure GitHub Releases has binaries named ${BINARY_NAME}-${OS}-${ARCH}"
    exit 1
fi

chmod +x "${TMP_DIR}/${BINARY_NAME}"

echo "Installing to /usr/local/bin (requires sudo)..."
if [ -w "/usr/local/bin" ]; then
    cp "${TMP_DIR}/${BINARY_NAME}" "/usr/local/bin/${BINARY_NAME}"
else
    sudo cp "${TMP_DIR}/${BINARY_NAME}" "/usr/local/bin/${BINARY_NAME}"
fi

echo "Success! ${BINARY_NAME} has been installed to /usr/local/bin"
echo ""

# Run setup
if command -v wiregate >/dev/null 2>&1; then
    echo "Running local setup..."
    wiregate setup
else
    echo "Warning: wiregate was installed but is not in your PATH yet."
    echo "Please restart your terminal or run: export PATH=\$PATH:/usr/local/bin"
fi
