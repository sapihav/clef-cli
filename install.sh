#!/bin/bash
set -e

# clef Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/sapihav/clef/main/install.sh | bash

REPO="sapihav/clef"
BINARY_NAME="clef"
DEFAULT_INSTALL_DIR="/usr/local/bin"
FALLBACK_INSTALL_DIR="${HOME}/.local/bin"
# Track whether user explicitly set INSTALL_DIR (no fallback if they did)
if [ -n "${INSTALL_DIR+x}" ]; then
    USER_SET_INSTALL_DIR=1
else
    USER_SET_INSTALL_DIR=0
    INSTALL_DIR="$DEFAULT_INSTALL_DIR"
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        darwin) OS="Darwin" ;;
        linux) OS="Linux" ;;
        *) error "Unsupported operating system: $OS" ;;
    esac
}

# Detect architecture
detect_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64) ARCH="x86_64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    if ! command -v jq &> /dev/null; then
        error "jq is required but not installed. Install with: brew install jq (macOS) or apt install jq (Linux)"
    fi
    LATEST_VERSION=$(curl -sS "https://api.github.com/repos/${REPO}/releases/latest" | jq -r '.tag_name // empty')
    if [ -z "$LATEST_VERSION" ]; then
        error "Failed to get latest version. Check your internet connection or GitHub API rate limits."
    fi
    if ! echo "$LATEST_VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then
        error "Unexpected version format: ${LATEST_VERSION}"
    fi
}

# Download and install
install() {
    detect_os
    detect_arch
    get_latest_version

    info "Installing clef ${LATEST_VERSION} for ${OS}/${ARCH}..."

    ARCHIVE_NAME="clef_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${ARCHIVE_NAME}"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/checksums.txt"

    info "Downloading from ${DOWNLOAD_URL}..."
    if ! curl -sSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE_NAME}"; then
        error "Failed to download. Release may not exist yet."
    fi

    info "Downloading checksums..."
    if ! curl -sSL "$CHECKSUM_URL" -o "${TMP_DIR}/checksums.txt"; then
        error "Failed to download checksums. Cannot verify integrity."
    fi

    # Verify checksum
    info "Verifying checksum..."
    cd "$TMP_DIR"
    EXPECTED=$(grep "${ARCHIVE_NAME}" checksums.txt | awk '{print $1}')
    if [ -z "$EXPECTED" ]; then
        error "No checksum found for ${ARCHIVE_NAME} in checksums.txt"
    fi
    if command -v sha256sum &> /dev/null; then
        ACTUAL=$(sha256sum "${ARCHIVE_NAME}" | awk '{print $1}')
    elif command -v shasum &> /dev/null; then
        ACTUAL=$(shasum -a 256 "${ARCHIVE_NAME}" | awk '{print $1}')
    else
        error "No SHA-256 tool found. Install coreutils (sha256sum) or shasum."
    fi
    if [ "$EXPECTED" != "$ACTUAL" ]; then
        error "Checksum mismatch!\n  Expected: ${EXPECTED}\n  Got:      ${ACTUAL}\nThe download may be corrupted or tampered with."
    fi
    info "Checksum verified."

    # Extract
    info "Extracting..."
    tar -xzf "$ARCHIVE_NAME"

    # Install — fall back to $HOME/.local/bin if default isn't writable
    # (common on macOS where /usr/local/bin needs sudo). Skip fallback if
    # user explicitly set INSTALL_DIR — respect their choice.
    if [ ! -d "$INSTALL_DIR" ] || [ ! -w "$INSTALL_DIR" ]; then
        if [ "$USER_SET_INSTALL_DIR" = "1" ]; then
            warn "${INSTALL_DIR} is not writable. Run with sudo or choose a writable path:"
            warn "  INSTALL_DIR=\$HOME/.local/bin bash install.sh"
            error "Cannot install without write access to ${INSTALL_DIR}"
        fi
        warn "${INSTALL_DIR} is not writable, falling back to ${FALLBACK_INSTALL_DIR}"
        INSTALL_DIR="$FALLBACK_INSTALL_DIR"
        mkdir -p "$INSTALL_DIR" || error "Failed to create ${INSTALL_DIR}"
    fi

    mv "$BINARY_NAME" "$INSTALL_DIR/"

    # Verify installation
    if command -v "$BINARY_NAME" &> /dev/null && [ "$(command -v "$BINARY_NAME")" = "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        info "Successfully installed clef to ${INSTALL_DIR}/${BINARY_NAME}"
        echo ""
        "$BINARY_NAME" --version
        echo ""
        info "Run 'clef --help' to get started"
    else
        info "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
        warn "${INSTALL_DIR} is not in your PATH — add it to use 'clef' directly:"
        echo ""
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
        echo ""
        echo "Or run directly: ${INSTALL_DIR}/${BINARY_NAME} --version"
    fi
}

# Run installer
install
