#!/bin/sh
set -e

# Switchboard Installation Script
# Install the latest version of Switchboard from GitHub releases

REPO="nickygerritsen/switchboard"
INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="switchboard"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
}

print_warning() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "Darwin"
            ;;
        Linux*)
            echo "Linux"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            print_error "This script supports macOS and Linux only."
            print_error "For Windows, use install.ps1 instead."
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "x86_64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            print_error "Switchboard supports x86_64 and arm64 only."
            exit 1
            ;;
    esac
}

# Get latest version from GitHub
get_latest_version() {
    print_info "Fetching latest version..."

    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi

    echo "$VERSION"
}

# Download file
download_file() {
    url="$1"
    output="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$output" "$url"
    else
        print_error "Neither curl nor wget found"
        exit 1
    fi
}

# Verify checksum
verify_checksum() {
    archive_file="$1"
    checksums_file="$2"

    print_info "Verifying checksum..."

    # Extract expected checksum for our file
    expected_checksum=$(grep "$(basename "$archive_file")" "$checksums_file" | awk '{print $1}')

    if [ -z "$expected_checksum" ]; then
        print_warning "Could not find checksum for $(basename "$archive_file")"
        return 1
    fi

    # Calculate actual checksum
    if command -v sha256sum >/dev/null 2>&1; then
        actual_checksum=$(sha256sum "$archive_file" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual_checksum=$(shasum -a 256 "$archive_file" | awk '{print $1}')
    else
        print_warning "Neither sha256sum nor shasum found, skipping checksum verification"
        return 0
    fi

    if [ "$expected_checksum" = "$actual_checksum" ]; then
        print_info "Checksum verified successfully"
        return 0
    else
        print_error "Checksum verification failed!"
        print_error "Expected: $expected_checksum"
        print_error "Got:      $actual_checksum"
        return 1
    fi
}

# Main installation
main() {
    print_info "Installing Switchboard..."

    # Parse command line arguments
    VERSION=""
    REGISTER=0
    while [ $# -gt 0 ]; do
        case "$1" in
            --version=*)
                VERSION="${1#*=}"
                shift
                ;;
            --register)
                REGISTER=1
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Usage: $0 [--version=VERSION] [--register]"
                exit 1
                ;;
        esac
    done

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)

    print_info "Detected OS: $OS"
    print_info "Detected Architecture: $ARCH"

    # Get version if not specified
    if [ -z "$VERSION" ]; then
        VERSION=$(get_latest_version)
    fi

    print_info "Installing version: $VERSION"

    # Build download URLs
    ARCHIVE_NAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE_NAME"
    CHECKSUMS_URL="https://github.com/$REPO/releases/download/$VERSION/checksums.txt"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download archive and checksums
    print_info "Downloading $ARCHIVE_NAME..."
    download_file "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE_NAME"

    print_info "Downloading checksums..."
    download_file "$CHECKSUMS_URL" "$TMP_DIR/checksums.txt"

    # Verify checksum
    if ! verify_checksum "$TMP_DIR/$ARCHIVE_NAME" "$TMP_DIR/checksums.txt"; then
        print_error "Installation aborted due to checksum verification failure"
        exit 1
    fi

    # Extract archive
    print_info "Extracting archive..."
    tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"

    # Determine installation directory
    if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" = "0" ]; then
        TARGET_DIR="$INSTALL_DIR"
        USE_SUDO=""
    else
        # Check if we can use sudo
        if command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
            TARGET_DIR="$INSTALL_DIR"
            USE_SUDO="sudo"
        else
            print_warning "$INSTALL_DIR is not writable and sudo is not available"
            print_info "Installing to $USER_INSTALL_DIR instead"
            TARGET_DIR="$USER_INSTALL_DIR"
            USE_SUDO=""

            # Create user install directory if it doesn't exist
            mkdir -p "$TARGET_DIR"
        fi
    fi

    # Install binary
    print_info "Installing to $TARGET_DIR..."
    if [ -n "$USE_SUDO" ]; then
        sudo install -m 755 "$TMP_DIR/$BINARY_NAME" "$TARGET_DIR/$BINARY_NAME"
    else
        install -m 755 "$TMP_DIR/$BINARY_NAME" "$TARGET_DIR/$BINARY_NAME"
    fi

    # Verify installation
    if ! command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_warning "$BINARY_NAME installed but not found in PATH"

        if [ "$TARGET_DIR" = "$USER_INSTALL_DIR" ]; then
            print_warning "You may need to add $USER_INSTALL_DIR to your PATH"
            print_info "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            print_info "  export PATH=\"\$PATH:$USER_INSTALL_DIR\""
        fi
    else
        print_info "Installation successful!"
        print_info "Installed: $TARGET_DIR/$BINARY_NAME"

        # Show version
        INSTALLED_VERSION=$("$BINARY_NAME" --version 2>&1 | head -n1)
        print_info "Version: $INSTALLED_VERSION"
    fi

    # Optionally register as default browser
    if [ $REGISTER -eq 1 ]; then
        print_info "Registering Switchboard as a browser..."
        if "$BINARY_NAME" register; then
            print_info "Registration successful!"
            print_info "You can now set Switchboard as your default browser in system settings"
        else
            print_warning "Registration failed. You can register manually later with: switchboard register"
        fi
    else
        print_info ""
        print_info "To register Switchboard as a browser, run: switchboard register"
        print_info "Then set it as your default browser in system settings"
    fi

    print_info ""
    print_info "Installation complete!"
}

# Run main function with all arguments
main "$@"
