# Git Hooks with Lefthook

This project uses [Lefthook](https://github.com/evilmartians/lefthook) to manage Git hooks for automated code quality checks and commit message formatting.

## Setup

### Prerequisites

1. **Lefthook**: Install lefthook
   ```bash
   # Using Homebrew (macOS)
   brew install lefthook
   
   # Using Go
   go install github.com/evilmartians/lefthook@latest
   ```

2. **golangci-lint**: Install Go linter (optional, for enhanced linting)
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Installation

After cloning the repository, install the hooks:

```bash
lefthook install
```

This will set up the following Git hooks:

## Hooks Overview

### Pre-commit Hook

Runs automatically before every commit with the following checks:

- **Format**: `gofmt -w .` - Automatically formats Go code
- **Vet**: `go vet ./...` - Static analysis for common Go issues  
- **Test**: `go test ./...` - Runs all unit tests
- **Build**: `go build` - Verifies the code compiles

All checks run in parallel for faster execution. If any check fails, the commit is blocked.

### Prepare-commit-msg Hook

Automatically prefixes commit messages based on branch name:

- `feat/new-feature` → `feat: your message`
- `fix/bug-description` → `fix: your message`
- `docs/documentation` → `docs: your message`
- `chore/maintenance` → `chore: your message`

Supported prefixes: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`

### Commit-msg Hook

Validates commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) format:

- **Format**: `<type>: <description>`
- **Types**: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert
- **Examples**:
  - ✅ `feat: add user authentication`
  - ✅ `fix: resolve login validation error`
  - ✅ `docs: update API documentation`
  - ❌ `added new feature` (missing type)
  - ❌ `Fix Bug` (should be lowercase)

## Configuration Files

### lefthook.yml

Main configuration file defining all hooks and their commands.

### .golangci.yml

Configuration for golangci-lint with project-specific rules and exclusions.

### scripts/

- `add-branch-prefix.sh`: Shell script for automatic commit message prefixing
- `validate-commit-msg.go`: Go program for commit message validation

## Usage

### Normal Development

Once installed, hooks run automatically:

1. **Make changes** to your code
2. **Stage changes**: `git add .`
3. **Commit**: `git commit -m "your message"`
   - Pre-commit checks run automatically
   - Commit message gets prefixed based on branch name
   - Commit message format is validated

### Manual Hook Execution

Test hooks without committing:

```bash
# Run all pre-commit hooks
lefthook run pre-commit

# Run specific hook
lefthook run pre-commit format
lefthook run pre-commit test
```

### Bypassing Hooks

In exceptional cases, you can skip hooks:

```bash
# Skip all hooks
git commit --no-verify -m "emergency fix"

# Skip specific hook (not recommended)
LEFTHOOK_EXCLUDE=test git commit -m "skip tests"
```

## Branch Naming Convention

To take advantage of automatic commit message prefixing, use this naming pattern:

```
<type>/<description>
```

Examples:
- `feat/user-authentication`
- `fix/login-validation-error`
- `docs/api-documentation-update`
- `chore/dependency-updates`
- `refactor/cleanup-validation-logic`

## Troubleshooting

### Hook Installation Issues

```bash
# Reinstall hooks
lefthook uninstall
lefthook install

# Verify installation
ls -la .git/hooks/
```

### Hook Execution Issues

```bash
# Check lefthook configuration
lefthook version
lefthook info

# Run with verbose output
LEFTHOOK_VERBOSE=1 lefthook run pre-commit
```

### Permission Issues

```bash
# Make scripts executable
chmod +x scripts/add-branch-prefix.sh
```

### Skipping Problematic Checks

If a specific check is causing issues, you can temporarily disable it:

```bash
# Skip linting
LEFTHOOK_EXCLUDE=lint git commit -m "temporary fix"

# Skip tests
LEFTHOOK_EXCLUDE=test git commit -m "work in progress"
```

## Customization

### Adding New Hooks

Edit `lefthook.yml` to add new commands or modify existing ones:

```yaml
pre-commit:
  commands:
    new-check:
      run: your-command-here
      glob: "*.go"
```

### Modifying Linter Rules

Edit `.golangci.yml` to customize linting rules:

```yaml
linters:
  enable:
    - your-preferred-linters
```

### Custom Commit Message Patterns

Modify `scripts/add-branch-prefix.sh` to support additional branch name patterns:

```bash
case "$BRANCH_NAME" in
    your-pattern/*)
        COMMIT_TYPE="your-type"
        ;;
esac
```

## Benefits

- **Consistency**: Automated code formatting and conventional commits
- **Quality**: Pre-commit testing and analysis prevent broken code
- **Efficiency**: Parallel execution minimizes hook execution time  
- **Standards**: Enforced coding standards and commit message format
- **Automation**: Reduces manual overhead and human error

## Integration with CI/CD

The same checks that run locally will also run in CI/CD pipelines, ensuring consistency across development environments.