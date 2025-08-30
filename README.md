# gh-topic-urls

A CLI tool to fetch GitHub Pull Request URLs for a specific branch and copy them to the clipboard in Markdown list format.

![result](https://github.com/user-attachments/assets/5e64d4e8-eec8-43fa-a1be-a654a158b66a)


## Features

- **Dynamic repository detection** - Automatically detects the current Git repository
- **Smart branch handling** - Uses current branch when no argument is provided
- **Branch validation** - Verifies branch existence before processing
- **User-friendly error messages** - Clear English error messages
- **Conditional clipboard copy** - Only copies to clipboard when PRs are found
- **Markdown formatting** - Formats URLs as Markdown list items
- **Timeout handling** - 30-second timeout for API requests

## Prerequisites

- [GitHub CLI](https://cli.github.com/) (`gh`) must be installed and authenticated
- `jq` command must be available for JSON processing

## Installation

```bash
gh extension install Yuki-Sakaguchi/gh-topic-urls
```

## Usage

```bash
# Use current branch (recommended)
gh topic-urls

# Specify a branch
gh topic-urls <branch-name>
```

### Examples

```bash
# Fetch PRs for the current branch
gh topic-urls

# Fetch PRs for the 'main' branch
gh topic-urls main

# Fetch PRs for a feature branch
gh topic-urls feature/new-feature
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

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access

## License

This project is licensed under the MIT License.
