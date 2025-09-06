# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-09-06

### Added
- Interactive branch selection with `--interactive` flag using promptui
- Shell auto-completion for branch names in bash/zsh/fish
- Branch listing functionality with `getAllBranches()` function
- Integration with GitHub CLI extension completion system
- Zero-setup shell completion (automatic when installed as GitHub CLI extension)
- Comprehensive test suite with unit and integration tests
- Branch name normalization and filtering logic

### Changed
- Optimized shell completion for GitHub CLI extensions
- Updated README with interactive selection and auto-completion documentation
- Enhanced error handling for git command failures
- Improved user experience with interactive prompts

### Technical
- Added `github.com/manifoldco/promptui v0.9.0` dependency
- Implemented TDD approach with proper test coverage
- Added comprehensive integration tests for branch selection functionality

## [1.0.0] - 2025-08-30

### Added
- CI/CD pipeline with automated semantic versioning
- Release branch strategy for controlled releases
- Cross-platform binary builds with GoReleaser
- Automated changelog generation
- Git hooks with lefthook for code quality
- GitHub CLI extension for fetching PR URLs by branch
- Dynamic repository detection from Git remote
- Smart branch handling with current branch fallback
- Branch validation before processing
- Clipboard integration for easy URL copying
- Markdown formatting for PR URLs
- Comprehensive error handling and user-friendly messages
- 30-second timeout for API requests

### Changed
- Development workflow now targets `release/next` branch
- Updated contribution guidelines in README

### Added
- GitHub CLI extension for fetching PR URLs by branch
- Dynamic repository detection from Git remote
- Smart branch handling with current branch fallback
- Branch validation before processing
- Clipboard integration for easy URL copying
- Markdown formatting for PR URLs
- Comprehensive error handling and user-friendly messages
- 30-second timeout for API requests