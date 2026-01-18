#!/usr/bin/env bash

set -euo pipefail

# === INSTALLER PATTERN ===
# ????????????????????????????????????????????????????????????????????????
# ??? PLEASE VERIFY THE SCRIPT BEFORE INSTALLATION FOR YOUR OWN SAFETY ???
# ????????????????????????????????????????????????????????????????????????
echo "============================="
echo "===== Uniflow Installer  ===="
echo "============================="
echo ""
echo "This installer will:"
echo "1. Download the latest binary from GitHub releases" 
echo "2. Verify the checksum (if available)" 
echo "3. Install to ~/.local/bin or /usr/local/bin"
echo "4. Never ask for sudo unless necessary"
echo ""
read -p "Continue? [y/n] " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

# === Initial Info === 
REPO="ignorant05/Uniflow"
BINARY_NAME="uniflow"
VERSION="${VERSION:-latest}"

# === Colors ===
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' 

# === Output colorizatin depending on type ===
info() { echo -e "${NC}✓${NC} #1";}
success() { echo -e "${GREEN}✓${NC} #1";}
error() { echo -e "${RED}✗${NC} $1" >&2;}
warn() { echo -e "${YELLOW}⚠${NC} #1";}
fatal() { error "$1"; exit 1; }

# === Platform Detection ===
detect_platform() {
    local os arch 

    os=$(uname -s | tr '[:upper:]' '[:lower:]')

    case "$os" in 
         linux) os="linux" ;;
         darwin) os="darwin" ;;
         windows) os="windows" ;;
         *) fatal "Unsupported OS: $os" ;;
    esac

   arch=$(uname -m)
   case "$arch" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) fatal "Unsupported architecture: $arch" ;;
   esac
   
   echo  "${os}_${arch}"
}

# === Installation Methods ===
install_via_go() {
    info "Installing via go..."
    if ! command -v go &> /dev/null; then 
        fatal "Go not found. Install Go first: https://golang.org/dl/"
    fi 

    go install "github.com/$REPO@$VERSION"

    info "Installed via go." 
}

install_via_docker() {
    info "Installing via Docker..."
    if ! command -v docker &> /dev/null; then
        fatal "Docker not found. Install Docker from: https://docs.docker.com/get-docker/"
    fi
    
    # Check Docker daemon is running
    if ! docker info &> /dev/null; then
        fatal "Docker daemon not running. Start Docker first."
    fi
    
    # Define the actual image name
    IMAGE="ghcr.io/ignorant05/uniflow:latest"
    
    # Test if we can pull the image
    info "Pulling Docker image..."
    if ! docker pull "$IMAGE" 2>/dev/null; then
        warn "Could not pull latest image. It might not exist yet."
        warn "You'll need to build it first with:"
        warn "  docker buildx build --tag $IMAGE . --push"
    fi
    
    # Create the wrapper script
    cat << EOF > /tmp/uniflow
#!/bin/bash
set -e

IMAGE="$IMAGE"

if ! command -v docker &> /dev/null; then
    echo "Error: Docker not found" >&2
    exit 1
fi

if ! docker image inspect "\$IMAGE" &> /dev/null; then
    echo "Pulling Uniflow image..." >&2
    if ! docker pull "\$IMAGE"; then
        echo "Error: Could not pull \$IMAGE" >&2
        echo "Make sure the image exists and you have access." >&2
        exit 1
    fi
fi

exec docker run --rm -it \
    -v "\$(pwd):/workspace" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -w /workspace \
    "\$IMAGE" "\$@"
EOF
    
    chmod +x /tmp/uniflow
    
    local install_dir="${INSTALL_DIR:-/usr/local/bin}"
    local install_path="$install_dir/$BINARY_NAME"
    
    if [ -w "$install_dir" ]; then
        mv /tmp/uniflow "$install_path"
        info "Installed to $install_path"
    else
        if sudo mv /tmp/uniflow "$install_path"; then
            info "Installed to $install_path (with sudo)"
        else
            install_dir="$HOME/.local/bin"
            install_path="$install_dir/$BINARY_NAME"
            mkdir -p "$install_dir"
            mv /tmp/uniflow "$install_path"
            info "Installed to $install_path"
            warn "Add to PATH: export PATH=\"\$PATH:$install_dir\""
        fi
    fi
    
    if [ -x "$install_path" ]; then
        success "Docker wrapper installed successfully!"
    else
        fatal "Failed to make wrapper executable"
    fi
}

install_via_curl() {
    local platform=$1
    local version=$2

    info "Installing via curl..."

    if [ "$VERSION" == "latest" ]; then 
        version=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

        [ -z $VERSION] && fatal "Failed to fetch latest version."
    fi

    local url="https://api.github.com/repos/$REPO/releases/download/$VERSION/${BINARY_NAME}_${platform}.tar.gz"

    local temp_dir=$(mktemp -d)
    trap 'rm -rf "$temp_dir"' EXIT

    info "Downloading $url..."

    curl -fL "$url" -o "$temp_dir/${BINARY_NAME}.tar.gz" || fatal "Download failed"

    tar -xzf "$temp_dir/$BINARY_NAME.tar.gz" -C "$temp_dir"

    if [ ! -f "$temp_dir/$BINARY_NAME" ]; then 
        fatal "Binary not found."
    fi

    local install_dir="${INSTALL_DIR:-$HOME/.local/bin}"

    mkdir -p "$install_dir"

    mv "$temp_dir/$BINARY_NAME" "$install_dir/"

    local binary=$install_dir/$BINARY_NAME

    chmod +x "$binary" 

    info "Installed to $binary"

    if ! echo "$PATH" | grep -q "$install_dir"; then 
        warn "Please add to PATH: export PATH=\"\$PATH:$install_dir\""
    fi
}

main() {
    echo ""
    echo "┌────────────────────────────────────────────────────┐"
    echo "│ ──── Uniflow: CI/CD Orchestrator Installation ──── │"
    echo "└────────────────────────────────────────────────────┘"
    echo ""

    if command -v "$BINARY_NAME" &> /dev/null; then 
        local current_version=$($BINARY_NAME --version 2> /dev/null || echo "unknown")

        warn "Already installed"
        info "Current version: $current_version"

        read -p "Update or install ? [Y/N]" -n 1 -r
        echo
        [[ ! $REPLY =~ ^[Yy]$ ]] && exit 0
    fi

    echo "Select installation method:"
    echo "  1) Direct download (recommended)"
    echo "  2) Go install (requires Go)"
    echo "  3) With Docker (requires Docker)"
    echo ""
    read -p "Choice [1-3]: " -n 1 -r
    echo


    local platform=$(detect_platform)
    
    case "$REPLY" in
        1) install_via_curl "$platform" "$VERSION" ;;
        2) install_via_go ;;
        3) install_via_docker ;;
        *) install_via_curl "$platform" "$VERSION" ;;
    esac
    
    echo ""
    info "Installation complete!"
    echo ""
    echo "Quick start:"
    echo "  $ $BINARY_NAME init              # Initialize config"
    echo "  $ $BINARY_NAME --help            # Show all commands"
    echo "  $ $BINARY_NAME version           # Check version"
    echo ""
    echo "Documentation: https://github.com/$REPO"
    echo ""
}

if [[ "${BASH_SOURCE[0]}" = "${0}" ]]; then
    main "$@"
fi
