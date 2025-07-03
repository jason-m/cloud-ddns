#!/usr/bin/env bash

# Cloud DDNS Multi-Platform Build Script
# Builds binaries for Linux, macOS, and FreeBSD on x86_64 and arm64

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Application name
APP_NAME="cloud-ddns"

# Build directory
BUILD_DIR="build"

# Version (can be overridden with BUILD_VERSION env var)
VERSION=${BUILD_VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}

echo -e "${BLUE}Building ${APP_NAME} v${VERSION}${NC}"
echo -e "${BLUE}================================${NC}"

# Create build directory
mkdir -p ${BUILD_DIR}

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf ${BUILD_DIR}/*

# Define target platforms
declare -a platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "freebsd/amd64"
    "freebsd/arm64"
)

# Build function
build_binary() {
    local goos=$1
    local goarch=$2
    
    # Convert Go OS names to our naming convention
    case $goos in
        "darwin")
            os_name="macos"
            ;;
        "linux")
            os_name="linux"
            ;;
        "freebsd")
            os_name="freebsd"
            ;;
        *)
            os_name=$goos
            ;;
    esac
    
    # Convert Go arch names to our naming convention
    case $goarch in
        "amd64")
            arch_name="x86_64"
            ;;
        "arm64")
            arch_name="arm64"
            ;;
        *)
            arch_name=$goarch
            ;;
    esac
    
    local output_name="${APP_NAME}-${os_name}_${arch_name}"
    local output_path="${BUILD_DIR}/${output_name}"
    
    echo -e "${YELLOW}Building for ${os_name}/${arch_name}...${NC}"
    
    # Set build environment
    export GOOS=$goos
    export GOARCH=$goarch
    export CGO_ENABLED=0
    
    # Build with ldflags for version info and smaller binary
    go build \
        -ldflags="-w -s -X main.version=${VERSION}" \
        -trimpath \
        -o "${output_path}" \
        .
    
    if [ $? -eq 0 ]; then
        # Get file size before compression
        size_before=$(du -h "${output_path}" | cut -f1)
        
        # Make executable
        chmod +x "${output_path}"
        
        # Gzip the binary
        gzip "${output_path}"
        local gzipped_path="${output_path}.gz"
        
        # Get file size after compression
        size_after=$(du -h "${gzipped_path}" | cut -f1)
        
        echo -e "${GREEN}✓ Built ${output_name}.gz (${size_before} → ${size_after})${NC}"
    else
        echo -e "${RED}✗ Failed to build ${output_name}${NC}"
        return 1
    fi
}

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

# Check if we're in a Go module
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Make sure you're in the project root directory.${NC}"
    exit 1
fi

# Download dependencies
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod download
go mod tidy

# Build for all platforms
echo -e "${YELLOW}Starting builds...${NC}"
for platform in "${platforms[@]}"; do
    IFS='/' read -r goos goarch <<< "$platform"
    build_binary "$goos" "$goarch"
done

# Reset environment
unset GOOS GOARCH CGO_ENABLED

echo -e "${BLUE}================================${NC}"
echo -e "${GREEN}Build completed!${NC}"
echo -e "${BLUE}Build artifacts:${NC}"

# List all built binaries with sizes
ls -lh ${BUILD_DIR}/ | grep -v "^total" | while read -r line; do
    echo -e "  ${GREEN}$(echo $line | awk '{print $9}')${NC} ($(echo $line | awk '{print $5}'))"
done

echo -e "${BLUE}================================${NC}"

# Optional: Create checksums
if command -v sha256sum &> /dev/null; then
    echo -e "${YELLOW}Generating checksums...${NC}"
    cd ${BUILD_DIR}
    sha256sum ${APP_NAME}-*.gz > checksums.txt
    cd ..
    echo -e "${GREEN}✓ Checksums saved to ${BUILD_DIR}/checksums.txt${NC}"
elif command -v shasum &> /dev/null; then
    echo -e "${YELLOW}Generating checksums...${NC}"
    cd ${BUILD_DIR}
    shasum -a 256 ${APP_NAME}-*.gz > checksums.txt
    cd ..
    echo -e "${GREEN}✓ Checksums saved to ${BUILD_DIR}/checksums.txt${NC}"
fi

echo -e "${GREEN}All builds completed successfully!${NC}"
