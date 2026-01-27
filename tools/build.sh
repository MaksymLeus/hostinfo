#!/usr/bin/env bash
set -euo pipefail
#
# build.sh — Build system for the Go application.
#
# Features:
#   ✔ Quick local build
#   ✔ Build all supported platforms
#   ✔ Build for a specific OS and/or architecture
#   ✔ Version auto-detection from Git
#   ✔ Version embedding into the binary
#   ✔ Clean build artifacts
#   ✔ Informative help message
#   ✔ Create SHA256 checksums: ./build.sh checksums
#

#========================================
# COLORS
#========================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'  
BLUE='\033[1;34m'

#========================================
# CONFIGURATION
#========================================
APP_NAME="hostinfo"
BUILD_DIR="bin"
MAIN_FILE="./cmd/server/hostinfo.go"

SUPPORTED_OSES=("linux" "darwin" "windows")
SUPPORTED_ARCHES=("amd64" "arm64")

VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")

#========================================
# HELP MESSAGE
#========================================

print_header() {
    echo -e "${BLUE}========================================================${NC}"
    echo -e "${BLUE}  HostInfo Build System — Version: ${VERSION}${NC}"
    echo -e "${BLUE}========================================================${NC}"
}

help_msg() {
cat <<EOF
Usage:
  ./build.sh [command] [options]

Commands:
  quick                  Build for the current system only (default)
  all                    Build for all supported platforms
  os <GOOS>              Build for a specific OS (builds all its architectures)
  arch <GOARCH>          Build for a specific architecture across all OSes
  target <GOOS> <GOARCH> Build for a specific OS + architecture
  list                   Show all supported OS/Architecture combinations
  clean                  Remove build outputs
  help                   Show help message

Examples:
  ./build.sh                      Quick local build
  ./build.sh all                  Build everything into ./bin/
  ./build.sh os linux             Build all linux binaries
  ./build.sh arch arm64           Build all arm64 binaries
  ./build.sh target darwin amd64  Build macOS Intel binary
  ./build.sh checksums            Generate SHA256 checksums for bin/ artifacts
  ./build.sh clean                Remove build artifacts
  ./build.sh list                 List all supported OS/ARCH combinations

Supported OS:   ${SUPPORTED_OSES[*]}
Supported Arch: ${SUPPORTED_ARCHES[*]}

Output:
  Binaries are stored in: ${BUILD_DIR}/

EOF
}


#========================================
# BUILD FUNCTION
#========================================
build_platform() {
    # Build a single platform binary
    local os="$1"
    local arch="$2"
    # if $3 is provided, use it as output, else use default naming
    local output="${BUILD_DIR}/${APP_NAME}-${os}-${arch}"
    [[ $# -gt 2 ]] && output="$3"

    # Create build directory if it doesn't exist
    mkdir -p "$BUILD_DIR"

    # Windows extension
    [[ "$os" == "windows" ]] && output+=".exe"

    echo -e "${YELLOW}→ Building for ${os}/${arch}${NC}"

    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
        go build -o "$output" \
        -ldflags="-s -w -X main.Version=${VERSION}" \
        "$MAIN_FILE"

    echo -e "${GREEN}✔ Built: ${output}${NC}"
}

#========================================
# FUNCTIONS
#========================================
list_platforms() {
    # List all supported OS/ARCH combinations
    echo -e "${GREEN}Supported platforms:${NC}"
    
    for os in "${SUPPORTED_OSES[@]}"; do
        for arch in "${SUPPORTED_ARCHES[@]}"; do
            echo "  - ${os}/${arch}"
        done
    done

}

build_all() {
    # Build for all supported OS/ARCH combinations
    echo -e "${CYAN}Building for all supported platforms...${NC}"

    for os in "${SUPPORTED_OSES[@]}"; do
        for arch in "${SUPPORTED_ARCHES[@]}"; do
            build_platform "$os" "$arch"
        done
    done

    echo -e "${GREEN}All builds completed successfully!${NC}"
    echo -e "${GREEN}Binaries are in the directory:${NC}" ${BUILD_DIR}/
    
    ls -lh "$BUILD_DIR"
}

build_for_os() {
    # Build for a specific OS across all architectures
    local target_os="$1"

    echo -e "${CYAN}Building for OS: ${target_os}${NC}"

    for os in "${SUPPORTED_OSES[@]}"; do
        for arch in "${SUPPORTED_ARCHES[@]}"; do
            [[ "$os" == "$target_os" ]] && build_platform "$os" "$arch"
        done
    done
}


build_for_arch() {
    # Build for a specific architecture across all OSes
    local target_arch="$1"

    echo -e "${BLUE}Building for architecture: ${target_arch}${NC}"

    for os in "${SUPPORTED_OSES[@]}"; do
        for arch in "${SUPPORTED_ARCHES[@]}"; do
            [[ "$arch" == "$target_arch" ]] && build_platform "$os" "$arch"
        done
    done
}

quick_build() {
    # Build for current OS/ARCH
    local output="./bin/${APP_NAME}"
    echo -e "${YELLOW}Quick local build...${NC}"
    build_platform "$(go env GOOS)" "$(go env GOARCH)" "$output"
    echo -e "${GREEN}✔ Built local binary: ${output}${NC}"
}

clean() {
    # Remove build outputs
    echo -e "${YELLOW}Cleaning...${NC}"
    rm -rf "$BUILD_DIR" "${APP_NAME}" "${APP_NAME}.exe" || true
    echo -e "${GREEN}✔ Clean complete.${NC}"
}

generate_checksums() {
    # Generate SHA256 checksums for all binaries in BUILD_DIR
    echo -e "${YELLOW}Generating SHA256 checksums...${NC}"
    cd "$BUILD_DIR"
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum ./* > checksums.txt
    else
        shasum -a 256 ./* > checksums.txt
    fi

    echo -e "${GREEN}✔ Checksums saved to ${BUILD_DIR}/checksums.txt${NC}"

}

#----------------------------------------
# VALIDATION FUNCTIONS
#----------------------------------------
validate_os() {
    # Validate OS and ARCH inputs 
    if [[ ! " ${SUPPORTED_OSES[*]} " =~ " $1 " ]]; then
        echo -e "${RED}Unsupported OS: $1${NC}"
        echo -e "${RED}Supported OSes: ${SUPPORTED_OSES[*]}${NC}"
        exit 1
    fi
}
validate_arch() {
    # Validate ARCH input
    if [[ ! " ${SUPPORTED_ARCHES[*]} " =~ " $1 " ]]; then
        echo -e "${RED}Unsupported Architecture: $1${NC}"
        echo -e "${RED}Supported Architectures: ${SUPPORTED_ARCHES[*]}${NC}"
        exit 1
    fi
}

simple_detection() {
  # Check for Go installation
  command -v go >/dev/null 2>&1 || { echo "Go not installed"; exit 1; }

  # Check if inside a git repository
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo -e "${YELLOW}Warning: not a git repo, using version=${VERSION}${NC}"
  fi

}

# ------------------------------------------------------------
# Main Logic
# ------------------------------------------------------------
print_header
simple_detection

cmd="${1:-quick}"
case "$cmd" in
    quick)
        quick_build
        ;;
    all)
        build_all
        ;;
    os)
        [[ $# -lt 2 ]] && { echo -e "${RED}Missing GOOS${NC}"; exit 1; }
        validate_os "$2"
        build_for_os "$2"
        ;;
    arch)
        [[ $# -lt 2 ]] && { echo -e "${RED}Missing GOARCH${NC}"; exit 1; }
        validate_arch "$2"
        build_for_arch "$2"
        ;;
    target)
        [[ $# -lt 3 ]] && { echo -e "${RED}Need GOOS & GOARCH${NC}"; exit 1; }
        validate_os "$2"
        validate_arch "$3"
        build_platform "$2" "$3"
        ;;
    checksums)
        generate_checksums
        ;;
    list)
        list_platforms
        ;;
    clean)
        clean
        ;;
    help|--help|-h)
        help_msg
        ;;
    *)
        echo -e "${RED}Unknown command: $cmd${NC}"
        help_msg
        exit 1
        ;;
esac
