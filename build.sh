#!/bin/bash

# Build script for hostinfo
# Builds binaries for Linux x64, Mac OS (Intel and Apple Silicon)
# Usage:
#   ./build.sh        - Build hostinfo binary in root directory for quick testing
#   ./build.sh all    - Build all platform binaries and store in bin/ directory

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

APP_NAME="hostinfo"
VERSION="1.0.0"
BUILD_DIR="bin"

echo -e "${GREEN}Building Host Info v${VERSION}${NC}"

# Function to build for a specific platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3

    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT" -ldflags="-s -w" ./cmd/server/hostinfo.go
    echo -e "${GREEN}Built: ${OUTPUT}${NC}"
}

# If no arguments or "quick" argument, build for current platform in root directory
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}Building quick test binary...${NC}"
    go build -o hostinfo ./cmd/server/hostinfo.go
    echo -e "${GREEN}Build complete! Run with: ./hostinfo${NC}"
    exit 0
fi

# Build all platforms
if [ "$1" == "all" ]; then
    # Create bin directory if it doesn't exist
    mkdir -p "$BUILD_DIR"

    # Build for different platforms
    build_platform "linux" "amd64" "${BUILD_DIR}/${APP_NAME}-linux-x64"
    build_platform "darwin" "amd64" "${BUILD_DIR}/${APP_NAME}-darwin-intel"
    build_platform "darwin" "arm64" "${BUILD_DIR}/${APP_NAME}-darwin-arm64"

    echo ""
    echo -e "${GREEN}All builds completed successfully!${NC}"
    echo -e "${GREEN}Binaries are in the ${BUILD_DIR}/ directory:${NC}"
    ls -lh ${BUILD_DIR}/
    exit 0
fi

echo -e "${RED}Invalid argument. Use './build.sh' for quick build or './build.sh all' for all platforms${NC}"
exit 1