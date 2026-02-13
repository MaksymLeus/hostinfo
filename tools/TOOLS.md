<!-- omit in toc -->
# Tools Descriptions

<!-- omit in toc -->
## Table of Contents
- [Project Preparation Script (prepare.sh)](#project-preparation-script-preparesh)
  - [Overview](#overview)
  - [Quick Setup](#quick-setup)
  - [🧰 What the Script Does](#-what-the-script-does)
  - [🔁 Safe to Re-run](#-safe-to-re-run)
---

## Project Preparation Script (prepare.sh)
Location: [`tools/prepare.sh`](./prepare.sh)

### Overview

The `prepare.sh` script bootstraps the development environment for the Go project.

It ensures:
- Go is installed

- Go modules are initialized and downloaded
  
- Frontend dependencies are installed (if applicable)

- .env file is created
  
- Git hooks are configured
  
- Docker presence is checked (optional)

The script is idempotent and safe to run multiple times.

### Quick Setup

```bash
# From project root:
bash tools/prepare.sh
# Or
chmod +x tools/prepare.sh
./tools/prepare.sh
```

### 🧰 What the Script Does

1. Environment Validation
   Checks for:
   - ✅ Go installation

   - ✅ Docker (optional)

   - ✅ Git repository

2. Go Setup
   ```bash
   # If go.mod does not exist:
   go mod init <project-name>
   # Then
   go mod tidy
   go mod download
   ```
3. Frontend Setup (Optional)
  
    If this file exists: `frontend/package.json`
     ```bash
     # The script runs:
     npm install
     ```
     Otherwise, frontend setup is skipped.

4. Environment File Setup

    If `.env` does not exist:  
      - Copies from `.env.example`
      - OR creates default `.env`:

        ```bash
        HOSTINFO_PORT=8080
        HOSTINFO_HOST=0.0.0.0
        HOSTINFO_DEBUG=false
        ```

5. Git Hooks Setup

    If inside a Git repository:
    - Sets:
        ```bash
        git config core.hooksPath tools/scripts/git_hooks
        ```
    - Makes hooks executable

    This enables centralized hook management.
    
    For a full description, see [`GIT_HOOKS`](./scripts/git_hooks/GIT_HOOKS.md).

### 🔁 Safe to Re-run
The script:
 
 - Does not overwrite existing `.env`
 - Does not reinitialize `go.mod` if exists
 - Does not override **hooksPath** if already configured

---