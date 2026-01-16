# Semantic Release Guide

This document explains how to use Semantic Release in this repository to automatically manage versioning, changelogs, and GitHub releases.


## Table of Contents
- [Semantic Release Guide](#semantic-release-guide)
  - [Table of Contents](#table-of-contents)
    - [Overview](#overview)
    - [How It Works](#how-it-works)
    - [Commit message examples:](#commit-message-examples)
    - [Configuration](#configuration)
    - [Key options:](#key-options)
    - [Workflow Usage](#workflow-usage)
    - [Inputs:](#inputs)
    - [Secrets:](#secrets)
    - [Custom Versioning](#custom-versioning)
    - [Troubleshooting](#troubleshooting)



### Overview

Semantic Release automates versioning and releases based on commit messages following Conventional Commits
.

- Version is automatically incremented according to commit type:

  - `feat`: Minor version bump
  - `fix`: Patch version bump
  - `BREAKING` CHANGE: Major version bump

- Changelog is automatically generated.
- GitHub release is created automatically with the generated notes.
- Optional npm publish (if Node.js project).


### How It Works

1. Analyze commits since the last release.
2. Determine the next semantic version.
3. Update `CHANGELOG.md`.
4. Commit changes (if configured).
5. Create a GitHub release with the new version.

### Commit message examples:
```
feat: add healthcheck endpoint
fix: resolve crash on startup
feat!: remove old API endpoint (BREAKING CHANGE)
```

### Configuration 

The repository uses `.releaserc.yaml `for configuration:
```yaml
branches:
  - main
  - "fs/pl/*"
  - "+([0-9])?(.{+([0-9]),x}).x"

plugins:
  - "@semantic-release/commit-analyzer"
  - "@semantic-release/release-notes-generator"
  - "@semantic-release/changelog"
  - "@semantic-release/git"
  - "@semantic-release/github"
```
### Key options:

- `branches`: Defines release branches. Tags are only created on these branches.

- `plugins`: Control which steps run (analyze commits, generate notes, update changelog, push changes, publish to GitHub).

Note: For a pure Go project, npm publishing steps can be skipped. Only GitHub release is needed.

### Workflow Usage

The GitHub Actions workflow is configured as a reusable workflow:
```yaml
jobs:
  release:
    uses: ./.github/workflows/template-semantic-release.yml
    with:
      dry_run: false
      semantic_version: 20
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

```
### Inputs:
| Input Name    | Description |
| -------- | ------- |
| `dry_run`  | `true` → simulate release without pushing tags. `false` → create actual release.    |
| `semantic_version` | Version of semantic-release to use (optional, default managed internally).     |

### Secrets:
| Secret Name | Description |
| -------- | ------- |
| `GITHUB_TOKEN`  | Required for GitHub API authentication and pushing tags.    |

### Custom Versioning

Semantic Release automatically determines the next version based on commit messages:

`fix`: → patch (`1.0.1`)

`feat:` → minor (`1.1.0`)

`BREAKING CHANGE` → major (`2.0.0`)

You can override versioning temporarily using `inputs.semantic_version` in the workflow.

### Troubleshooting

1) Permission errors pushing tags:

   - Ensure `GITHUB_TOKEN` is set in workflow secrets.

   - Workflow must run on allowed branches only.

2) Invalid release branch error:

   - Ensure your branch is listed in `branches` in `.releaserc.yaml`.

3) Dry run vs actual run:

   - Dry run does not push or tag. Set `dry_run: false` to perform a real release.