# Git Hooks

Git hooks are scripts that run automatically at certain points in your Git workflow. They help enforce code quality, commit message conventions, and other automation tasks.

In **hostinfo**, we use two main Git hooks:

- `pre-commit` → runs before a commit is finalized
- `commit-msg` → validates commit messages

These hooks are stored in [`scripts/hooks`](../scripts/hooks/).

## 1. Pre-Commit Hook

Ensures code meets quality standards before committing.

**Purpose**

- Run Go formatting (gofmt)
- Run Go static analysis (go vet, staticcheck)
- Optionally, build the project or test Dockerfile
- Prevent commits if checks fail

**Script path**: [`scripts/hooks/pre-commit`](../scripts/hooks/pre-commit)

## 2. Commit-Message Hook

Enforces **Conventional Commits**.

**Purpose**

- Ensure commit messages follow a predefined convention
- Support automation tools like semantic-release
- Prevent ambiguous commit messages

**Script path**: [`scripts/hooks/commit-msg`](../scripts/hooks/commit-msg)

## 3. Installing Git Hooks

### 3.1 Manual Installation
```sh
chmod +x ./tools/setup-hooks.sh
chmod +x ./scripts/hooks/pre-commit
chmod +x ./scripts/hooks/commit-msg
```
### 3.2 Using Makefile (optional)
```makefile
install-hooks:
  chmod +x ./tools/setup-hooks.sh
  chmod +x ./scripts/hooks/*
	chmod +x .git/hooks/*
  ./tools/setup-hooks.sh
```
```bash
make install-hooks
```
## 4. Workflow Diagram

```text
┌───────────────┐
│ Developer     │
│ runs `git add`│
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ pre-commit    │
│ (gofmt, vet,  │
│ staticcheck)  │
└───────┬───────┘
        │ pass
        ▼
┌───────────────┐
│ commit-msg    │
│ (Conventional │
│ Commits check)│
└───────┬───────┘
        │ pass
        ▼
┌───────────────┐
│ Git Commit    │
│ stored locally│
└───────┬───────┘
        │ push
        ▼
┌───────────────┐
│ CI/CD Workflow│
│ - semantic-   │
│   release     │
│ - Docker push │
└───────────────┘
```
- **pre-commit** → prevents bad code from being committed
- **commit-msg** → ensures proper commit messages for automated releases
- **CI/CD** → triggers semantic-release, updates CHANGELOG, builds Docker image

## 5. Best Practices

- Keep hooks fast (avoid heavy tasks in pre-commit)
- Always make scripts executable (chmod +x)
- Include hooks in repo for team onboarding
- Update hooks when conventions or CI/CD changes

## 6. Benefits

- ✅ Consistent commit messages → automatic versioning
- ✅ Code quality enforced before commits
- ✅ Prevents broken code from entering repository
- ✅ Smooth onboarding for new developers

## 7. References

- [Git Hooks Documentation](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)