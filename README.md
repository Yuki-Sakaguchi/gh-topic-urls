# gh-topic-urls

A CLI tool to fetch GitHub Pull Request URLs for a specific branch and copy them to the clipboard in Markdown list format.

## Features

- Fetch all Pull Requests for a given branch using GitHub CLI
- Format URLs as Markdown list items
- Automatically copy results to clipboard
- Support for timeout handling (30 seconds)

## Prerequisites

- [GitHub CLI](https://cli.github.com/) (`gh`) must be installed and authenticated
- `jq` command must be available for JSON processing

## Installation

```bash
go install github.com/Yuki-Sakaguchi/gh-topic-urls@latest
```

Or build from source:

```bash
git clone https://github.com/Yuki-Sakaguchi/gh-topic-urls.git
cd gh-topic-urls
go build -o gh-topic-urls
```

## Usage

```bash
gh-topic-urls <branch-name>
```

### Example

```bash
# Fetch PRs for the 'main' branch
gh-topic-urls main

# Fetch PRs for a feature branch
gh-topic-urls feature/new-feature
```

The tool will:
1. Query GitHub API for all Pull Requests targeting the specified branch
2. Format the URLs as a Markdown list
3. Display the results in the terminal
4. Copy the formatted list to your clipboard

### Output Format

```markdown
- https://github.com/yesodco/yesod/pull/123
- https://github.com/yesodco/yesod/pull/124
- https://github.com/yesodco/yesod/pull/125
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access

## License

This project is licensed under the MIT License.
