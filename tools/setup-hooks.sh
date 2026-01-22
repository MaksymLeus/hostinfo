#!/usr/bin/env bash

git config core.hooksPath scripts/hooks
echo "✔ Git hooks installed via core.hooksPath"
echo "Current hooksPath: $(git config core.hooksPath)"

