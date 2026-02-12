#!/usr/bin/env bash
set -Eeuo pipefail

#############################################
# GLOBAL CONFIG
#############################################

ENV_FILE=".env"
ENV_TEMPLATE=".env.example"
HOOKS_PATH="tools/scripts/git_hooks"
FRONTEND_DIR="frontend"

#############################################
# LOGGING
#############################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info()    { echo -e "${BLUE}ℹ $1${NC}"; }
log_success() { echo -e "${GREEN}✔ $1${NC}"; }
log_warn()    { echo -e "${YELLOW}⚠ $1${NC}"; }
log_error()   { echo -e "${RED}✖ $1${NC}"; }

#############################################
# ERROR HANDLING
#############################################

trap 'log_error "Script failed at line $LINENO"; exit 1' ERR

#############################################
# PATH DISCOVERY (cross-platform safe)
#############################################

SCRIPTPATH=$0
if [ ! -e "$SCRIPTPATH" ]; then
  case $SCRIPTPATH in
    (*/*) exit 1;;
    (*) SCRIPTPATH=$(command -v -- "$SCRIPTPATH") || exit;;
  esac
fi

dir=$(
  cd -P -- "$(dirname -- "$SCRIPTPATH")" && pwd -P
) || exit

SCRIPTPATH=$dir/$(basename -- "$SCRIPTPATH") || exit
PROJECT_HOME="$(dirname "$(dirname "$SCRIPTPATH")")"

cd "$PROJECT_HOME"

#############################################
# HEADER
#############################################

echo "-----------------------------------------"
log_info "🚀 Preparing Go Project"
log_info "Project root: $PROJECT_HOME"
echo "-----------------------------------------"

#############################################
# VERSION UTILS
#############################################

normalize_version() {
  # Extract first x.y.z pattern from version string
  echo "$1" | grep -Eo '[0-9]+(\.[0-9]+)+' | head -n1
}

version_ge() {
  # returns 0 if $1 >= $2
  # usage: version_ge "1.20.3" "1.18.0"
  [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

#############################################
# COMMAND CHECKS WITH VERSION SUPPORT
#############################################
require_command() {
  local cmd="$1"
  local message="$2"
  local min_version="${3:-}"
  local version_flag="${4:---version}"

  if ! command -v "$cmd" >/dev/null 2>&1; then
    log_error "$message"
    exit 1
  fi

  if [ -n "$min_version" ]; then
    local raw_version

    if [ "$version_flag" = "version" ]; then
      raw_version=$("$cmd" version 2>/dev/null | head -n1)
    else
      raw_version=$("$cmd" "$version_flag" 2>/dev/null | head -n1)
    fi

    local current_version
    current_version=$(normalize_version "$raw_version")

    if [ -z "$current_version" ]; then
      log_error "Unable to detect version for $cmd"
      exit 1
    fi

    if version_ge "$current_version" "$min_version"; then
      log_success "$cmd version $current_version (>= $min_version)"
    else
      log_error "$cmd version $current_version is less than required $min_version"
      exit 1
    fi
  else
    log_success "$cmd detected"
  fi
}

check_optional_command() {
  local cmd="$1"
  local label="$2"
  local min_version="${3:-}"
  local version_flag="${4:---version}"

  if ! command -v "$cmd" >/dev/null 2>&1; then
    log_warn "$label not found (optional)"
    return
  fi

  local raw_version

  if [ "$version_flag" = "version" ]; then
    raw_version=$("$cmd" version 2>/dev/null | head -n1)
  else
    raw_version=$("$cmd" "$version_flag" 2>/dev/null | head -n1)
  fi

  local current_version
  current_version=$(normalize_version "$raw_version")

  if [ -n "$min_version" ] && [ -n "$current_version" ]; then
    if version_ge "$current_version" "$min_version"; then
      log_success "$label version $current_version (>= $min_version)"
    else
      log_warn "$label version $current_version is less than recommended $min_version"
    fi
  else
    log_success "$label detected"
  fi
}

check_requirements() {
  require_command go "Go is required. Install from https://go.dev/dl/" "1.24.0" "version"

  check_optional_command npm "npm" "11.6.0" "-v"
  check_optional_command node "Node.js" "24.13.0" "-v"
  check_optional_command docker "Docker"
}

#############################################
# GO SETUP
#############################################

setup_go_modules() {
  if [ ! -f "$PROJECT_HOME/go.mod" ]; then
    log_warn "🟡 go.mod not found — initializing module"
    go mod init "$(basename "$PROJECT_HOME")"
  fi
  log_info "📦 Tidying Go modules..."
  go mod tidy

  log_info "📦 Downloading Go module dependencies..."
  go mod download

  log_success "Go modules ready"
}

#############################################
# FRONTEND SETUP
#############################################
setup_frontend() {
  if [ -f "$PROJECT_HOME/$FRONTEND_DIR/package.json" ]; then
    require_command npm "npm required for frontend setup"

    log_info "📦 Installing frontend dependencies..."
    (cd "$PROJECT_HOME/$FRONTEND_DIR" && npm install --silent)
    log_success "Frontend dependencies installed"
  else
    log_warn "No frontend/package.json found — skipping frontend setup"
  fi
}

#############################################
# ENV FILE SETUP
#############################################

setup_env() {
  if [ -f "$PROJECT_HOME/$ENV_FILE" ]; then
    log_success "$ENV_FILE already exists"
    return
  fi

  if [ -f "$PROJECT_HOME/$ENV_TEMPLATE" ]; then
    log_info "📄 Creating $ENV_FILE from $ENV_TEMPLATE"
    cp "$PROJECT_HOME/$ENV_TEMPLATE" "$PROJECT_HOME/$ENV_FILE"
  else
    log_info "📄 Creating default $ENV_FILE"
    cat <<EOF > "$PROJECT_HOME/$ENV_FILE"
HOSTINFO_PORT=8080
HOSTINFO_HOST=0.0.0.0
HOSTINFO_DEBUG=false
EOF
  fi

  log_success "$ENV_FILE created"
}

#############################################
# GIT HOOKS SETUP
#############################################

setup_git_hooks() {
  # Check git binary
  if ! command -v git >/dev/null 2>&1; then
    log_warn "Git not installed — skipping hooks"
    return
  fi

  # Check if we are in a git repo
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    log_warn "Not a Git repository — skipping hooks"
    return
  fi

  # Ensure hooks directory exists
  if [ ! -d "$PROJECT_HOME/$HOOKS_PATH" ]; then
    log_warn "Hooks directory: not found — skipping"
    return
  fi

  CURRENT_HOOKS_PATH=$(git config --get core.hooksPath || echo "")

  if [ "$CURRENT_HOOKS_PATH" != "$HOOKS_PATH" ]; then
    log_info "🔗 Configuring Git hooks path"

    if compgen -G "$PROJECT_HOME/$HOOKS_PATH/*" > /dev/null; then
      chmod +x "$PROJECT_HOME/$HOOKS_PATH"/*
    fi

    git config core.hooksPath "$PROJECT_HOME/$HOOKS_PATH"
    log_success "Git hooks configured"
  else
    log_success "Git hooks already configured"
  fi
}

#############################################
# POST CHECKS
#############################################

print_next_steps() {
  echo "-----------------------------------------"
  log_success "🎉 Preparation complete!"

  if [[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]]; then
    log_warn "Consider adding Go bin to PATH:"
    echo "  export PATH=\"\$PATH:\$(go env GOPATH)/bin\""
  fi

  echo "-----------------------------------------"
}

#############################################
# MAIN
#############################################

main() {
  check_requirements
  
  setup_go_modules
  setup_frontend
  setup_env
  setup_git_hooks

  print_next_steps
}

main "$@"