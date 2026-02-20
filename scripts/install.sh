#!/usr/bin/env bash
# Go Platform Template - Universal Installer
# Usage: curl -fsSL https://install.example.com/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="murtadanazar/go-api-template"
BINARY_NAME="go-platform"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${VERSION:-latest}"

# Detect OS and architecture
detect_os_arch() {
    local os arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) os="windows" ;;
        *)          os="unknown" ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        armv7)          arch="armv7" ;;
        i386)           arch="386" ;;
        *)              arch="unknown" ;;
    esac
    
    echo "${os}:${arch}"
}

# Main installation
main() {
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "${BLUE}  Go Platform Template - Installer${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo ""
    
    # Detect platform
    local platform
    platform=$(detect_os_arch)
    local os arch
    IFS=':' read -r os arch <<< "$platform"
    
    if [ "$os" = "unknown" ] || [ "$arch" = "unknown" ]; then
        echo -e "${RED}✗ Unsupported OS/Architecture: $(uname -s) $(uname -m)${NC}"
        exit 1
    fi
    
    # Determine binary name
    local binary_name="$BINARY_NAME"
    [ "$os" = "windows" ] && binary_name="${binary_name}.exe"
    
    # Get download URL
    local download_url
    if [ "$VERSION" = "latest" ]; then
        # Get latest release from GitHub API
        echo -e "${BLUE}→${NC} Fetching latest release..."
        local release_json
        release_json=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest")
        VERSION=$(echo "$release_json" | grep '"tag_name"' | head -1 | cut -d'"' -f4 | sed 's/^v//')
        
        if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
            echo -e "${RED}✗ Failed to fetch latest release${NC}"
            exit 1
        fi
    fi
    
    download_url="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}-v${VERSION}-${os}-${arch}"
    [ "$os" = "windows" ] && download_url="${download_url}.exe"
    
    echo -e "${BLUE}→${NC} Platform: ${YELLOW}${os}/${arch}${NC}"
    echo -e "${BLUE}→${NC} Version: ${YELLOW}v${VERSION}${NC}"
    echo -e "${BLUE}→${NC} Download URL: ${YELLOW}${download_url}${NC}"
    echo ""
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Download binary
    local temp_file
    temp_file=$(mktemp)
    trap "rm -f $temp_file" EXIT
    
    echo -e "${BLUE}→${NC} Downloading..."
    if ! curl -fsSL "$download_url" -o "$temp_file"; then
        echo -e "${RED}✗ Download failed${NC}"
        exit 1
    fi
    
    # Make executable
    chmod +x "$temp_file"
    
    # Move to install directory
    local install_path="${INSTALL_DIR}/${binary_name}"
    if ! mv "$temp_file" "$install_path"; then
        echo -e "${RED}✗ Failed to install to ${install_path}${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Installation successful${NC}"
    echo ""
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "${BLUE}  Next Steps${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo ""
    
    # Check if install dir is in PATH
    if [[ ":$PATH:" == *":${INSTALL_DIR}:"* ]]; then
        echo -e "${GREEN}✓ ${INSTALL_DIR} is in your PATH${NC}"
        echo ""
        echo "  Start using:"
        echo -e "    ${YELLOW}${binary_name}${NC}"
    else
        echo -e "${YELLOW}⚠ ${INSTALL_DIR} is NOT in your PATH${NC}"
        echo ""
        echo "  Add to your shell profile:"
        echo -e "    ${YELLOW}export PATH=\"${INSTALL_DIR}:\$PATH\"${NC}"
        echo ""
        echo "  Or run directly:"
        echo -e "    ${YELLOW}${install_path}${NC}"
    fi
    
    echo ""
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "${BLUE}  Documentation${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo ""
    echo "  Website: https://github.com/${REPO}"
    echo "  Issues:  https://github.com/${REPO}/issues"
    echo ""
}

main "$@"
