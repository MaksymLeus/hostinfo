#!/usr/bin/env bash
set -euo pipefail

#############################################
# CONFIG
#############################################

ENV_FILE=".env"
ENV_TEMPLATE=".env.example"

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

cd $PROJECT_HOME

#############################################
# HEADER
#############################################

echo "-----------------------------------------"
echo "🚀 Preparing Go project"
echo "Project root: $PROJECT_HOME"
echo "-----------------------------------------"

#############################################
# CHECKS
#############################################

check_go() {
  if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go not found. Install Go first: https://go.dev/dl/"
    exit 1
  fi
  echo "✔ Go version: $(go version)"
}

check_docker() {
  if command -v docker >/dev/null 2>&1; then
    echo "✔ Docker detected: $(docker --version)"
  else
    echo "⚠ Docker not found (optional). Install if you plan to containerize."
  fi
}

check_git() {

  if [ -d "$PROJECT_HOME/.git" ]; then
    echo "✔ Git repository detected"
  else
    echo "⚠ Not a Git repository (skip hooks setup)"
  fi
}

#############################################
# GO MODULES & DEPENDENCIES
#############################################

setup_go_modules() {
  if [ ! -f "$PROJECT_HOME/go.mod" ]; then
    echo "🟡 go.mod not found, initializing..."
    go mod init "$(basename "$PROJECT_HOME")"
  fi

  echo "📦 Tidying Go modules..."
  go mod tidy
}

#############################################
# ENV FILE SETUP
#############################################

setup_env() {
  if [ -f "$PROJECT_HOME/$ENV_FILE" ]; then
    echo "✔ $ENV_FILE already exists"
  else
    if [ -f "$ENV_TEMPLATE" ]; then
      echo "📄 Creating $ENV_FILE from $ENV_TEMPLATE"
      cp "$ENV_TEMPLATE" "$ENV_FILE"
    else
      echo "📄 Creating default $ENV_FILE"
      cat <<EOF > "$ENV_FILE"
HOSTINFO_PORT=8080
HOSTINFO_ADDR=0.0.0.0
HOSTINFO_DEBUG=false
EOF
    fi
    echo "✔ env file created"
  fi
}

#############################################
# PRE-COMMIT HOOKS
#############################################


setup_git_hooks() {
  # Check git binary
  if ! command -v git >/dev/null 2>&1; then
    echo "⚠ git not found — skipping hook setup"
    return
  fi

  # Check if we are in a git repo
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "⚠ not a git repository — skipping hook setup"
    return
  fi

  # Ensure hooks directory exists
  if [ ! -d "scripts/hooks" ]; then
    echo "⚠ scripts/hooks directory not found — skipping hook setup"
    return
  fi

  # Check current hooksPath
  CURRENT_HOOKS_PATH=$(git config --get core.hooksPath || echo "")

  
  # Set hooksPath only if different
  if [ "$CURRENT_HOOKS_PATH" != "scripts/hooks" ]; then
    echo "🔗 configuring core.hooksPath to scripts/hooks"
    chmod +x scripts/hooks/*
    git config core.hooksPath scripts/hooks
  else
    echo "> core.hooksPath already set to scripts/hooks"
  fi

  echo "✔ Git hooks installed via core.hooksPath"
  echo "> Current hooksPath: $(git config core.hooksPath)"
}


#############################################
# MAIN
#############################################

check_go
check_docker
check_git
setup_go_modules
setup_env
setup_git_hooks

echo "-----------------------------------------"
echo "🎉 Preparation complete!"
echo "Next steps:"
[[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]] && echo "  ➜ export PATH=\"\$PATH:\$(go env GOPATH)/bin\""
echo "  ➜ For more information, check the documentation in folder docs/*"
echo "-----------------------------------------"
