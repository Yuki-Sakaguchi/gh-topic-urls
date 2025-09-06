# gh-topic-urls

A CLI tool to fetch GitHub Pull Request URLs for a specific branch and copy them to the clipboard in Markdown list format.

![result](https://github.com/user-attachments/assets/5e64d4e8-eec8-43fa-a1be-a654a158b66a)


## Features

- **Dynamic repository detection** - Automatically detects the current Git repository
- **Smart branch handling** - Uses current branch when no argument is provided
- **Interactive branch selection** - Select branches with an intuitive UI using `--interactive` flag
- **Shell auto-completion** - Tab completion for branch names in bash/zsh/fish
- **Branch validation** - Verifies branch existence before processing
- **User-friendly error messages** - Clear English error messages
- **Conditional clipboard copy** - Only copies to clipboard when PRs are found
- **Markdown formatting** - Formats URLs as Markdown list items
- **Timeout handling** - 30-second timeout for API requests

## Prerequisites

- [GitHub CLI](https://cli.github.com/) (`gh`) must be installed and authenticated
- `jq` command must be available for JSON processing

## Installation

### Using GitHub CLI (Recommended)
```bash
gh extension install Yuki-Sakaguchi/gh-topic-urls
```

### Download Binary
Download the latest binary from [Releases](https://github.com/Yuki-Sakaguchi/gh-topic-urls/releases) page.

### Build from Source
```bash
git clone https://github.com/Yuki-Sakaguchi/gh-topic-urls.git
cd gh-topic-urls
go build -o gh-topic-urls .
```

## Usage

```bash
# Use current branch (recommended)
gh topic-urls

# Specify a branch
gh topic-urls <branch-name>

# Interactive branch selection
gh topic-urls --interactive
gh topic-urls -i
```

### Examples

```bash
# Fetch PRs for the current branch
gh topic-urls

# Fetch PRs for the 'main' branch
gh topic-urls main

# Fetch PRs for a feature branch
gh topic-urls feature/new-feature

# Interactive branch selection (NEW!)
gh topic-urls -i
```

## Shell Auto-completion Setup

Enable tab completion for branch names:

### Bash
```bash
# Add to ~/.bashrc
eval "$(gh topic-urls completion bash)"

# Or generate completion script
gh topic-urls completion bash > /usr/local/share/bash-completion/completions/gh-topic-urls
```

### Zsh
```bash
# Add to ~/.zshrc
eval "$(gh topic-urls completion zsh)"

# Or for oh-my-zsh
gh topic-urls completion zsh > "${fpath[1]}/_gh-topic-urls"
```

### Fish
```bash
gh topic-urls completion fish | source

# Or save to file
gh topic-urls completion fish > ~/.config/fish/completions/gh-topic-urls.fish
```

### How it works

1. **Auto-detects repository** from your current Git remote
2. **Validates branch existence** (when specified)
3. **Queries GitHub API** for all Pull Requests targeting the branch
4. **Formats URLs** as a Markdown list
5. **Displays results** in the terminal
6. **Copies to clipboard** (only if PRs are found)

### Output Examples

**When PRs are found:**
```
Using current branch: main
- https://github.com/your-org/your-repo/pull/123
- https://github.com/your-org/your-repo/pull/124
✨ Copied to clipboard
```

**When no PRs exist:**
```
Target branch: feature/empty-branch
No pull requests found for branch 'feature/empty-branch'
```

**When branch doesn't exist:**
```
Target branch: nonexistent
branch 'nonexistent' does not exist
```

## Development

### CI/CD Pipeline

This project uses a comprehensive CI/CD pipeline with automated semantic versioning and release management.

#### Branch Strategy
- **Main branch**: Always contains released code
- **Release branch** (`release/next`): Accumulates unreleased changes
- **Feature branches**: Target `release/next` for normal development
- **Hotfix branches**: Target `main` for emergency fixes

#### Automated Workflows
- **PR Validation**: Runs tests, linting, and build verification on PRs to `release/next`
- **Release Preparation**: Weekly automated or manual release preparation
- **Automated Release**: Builds and releases when release PRs are merged to main
- **Hotfix Support**: Emergency release process for critical fixes

#### Git Hooks
This project includes Git hooks powered by [Lefthook](https://github.com/evilmartians/lefthook):
- Pre-commit: Code formatting, linting, testing, and build verification
- Commit message: Automatic prefixing and validation
- See [HOOKS.md](docs/HOOKS.md) for detailed setup instructions

### Contributing

1. Create feature branch from `release/next`
2. Make changes following conventional commit format
3. Submit PR targeting `release/next`
4. CI will validate your changes automatically

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access

## License

This project is licensed under the MIT License.
